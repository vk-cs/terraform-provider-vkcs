package kms_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"
)

func TestKMSKeyDataSourceRead(t *testing.T) {
	for _, tc := range []struct {
		name            string
		path            string
		deletionAllowed bool
	}{
		{name: "test-key", path: "/kms/user/v1/transit/keys/test-key", deletionAllowed: true},
		{name: "key /?#%", path: "/kms/user/v1/transit/keys/key%20%2F%3F%23%25", deletionAllowed: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := fmt.Sprintf(`{"data":{"name":%q,"type":"aes256-gcm96","deletion_allowed":%t}}`, tc.name, tc.deletionAllowed)
			config := newKeyDataSourceTestConfig(t, tc.name, tc.path, http.StatusOK, body)
			ds := kms.NewKeyDataSource()
			var schema datasource.SchemaResponse
			ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{
				kms.KeyFieldName: tc.name, kms.KeyFieldType: "stale-type", kms.KeyFieldDeletionAllowed: !tc.deletionAllowed,
			})
			require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "stale-id").HasError())

			ds.(datasource.DataSourceWithConfigure).Configure(context.Background(), datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
			response := datasource.ReadResponse{State: d}
			ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config{Schema: d.Schema, Raw: d.Raw}}, &response)
			d = response.State
			diags := response.Diagnostics
			require.False(t, diags.HasError(), "%v", diags)
			var attribute1 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute1).HasError())
			assert.Equal(t, tc.name, attribute1.ValueString())
			var attribute2 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeyFieldName), &attribute2).HasError())
			assert.Equal(t, tc.name, attribute2.ValueString())
			var attribute3 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeyFieldType), &attribute3).HasError())
			assert.Equal(t, "aes256-gcm96", attribute3.ValueString())
			var attribute4 types.Bool
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeyFieldDeletionAllowed), &attribute4).HasError())
			assert.Equal(t, tc.deletionAllowed, attribute4.ValueBool())
		})
	}
}

func TestKMSKeyDataSourceRead_errors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{name: "forbidden", status: http.StatusForbidden, body: `{"errors":["permission denied"]}`},
		{name: "invalid JSON", status: http.StatusOK, body: `invalid`},
		{name: "invalid field type", status: http.StatusOK, body: `{"data":{"deletion_allowed":"yes"}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newKeyDataSourceTestConfig(t, "test-key", "/kms/user/v1/transit/keys/test-key", tc.status, tc.body)
			ds := kms.NewKeyDataSource()
			var schema datasource.SchemaResponse
			ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.KeyFieldName: "test-key"})

			ds.(datasource.DataSourceWithConfigure).Configure(context.Background(), datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
			response := datasource.ReadResponse{State: d}
			ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config{Schema: d.Schema, Raw: d.Raw}}, &response)
			d = response.State
			diags := response.Diagnostics
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary(), "Error getting key")
			var attribute5 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute5).HasError())
			assert.Empty(t, attribute5.ValueString())
		})
	}
}

func TestKMSKeyDataSourceRead_clientError(t *testing.T) {
	ds := kms.NewKeyDataSource()
	var schema datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.KeyFieldName: "test-key"})
	config := &secretDataSourceTestConfig{err: errors.New("client unavailable")}

	ds.(datasource.DataSourceWithConfigure).Configure(context.Background(), datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
	response := datasource.ReadResponse{State: d}
	ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config{Schema: d.Schema, Raw: d.Raw}}, &response)
	d = response.State
	diags := response.Diagnostics
	require.True(t, diags.HasError())
	require.Len(t, diags, 1)
	assert.Equal(t, "Error creating VKCS KMS client", diags[0].Summary())
	assert.Equal(t, "client unavailable", diags[0].Detail())
	var attribute6 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute6).HasError())
	assert.Empty(t, attribute6.ValueString())
}

func newKeyDataSourceTestConfig(t *testing.T, name, path string, status int, body string) *secretDataSourceTestConfig {
	t.Helper()
	calls := 0
	t.Cleanup(func() { assert.Equal(t, 1, calls) })
	provider := &gophercloud.ProviderClient{}
	provider.HTTPClient.Transport = secretDataSourceTestTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, path, r.URL.EscapedPath())
		assert.Equal(t, "/kms/user/v1/transit/keys/"+name, r.URL.Path)
		assert.Empty(t, r.URL.RawQuery)
		assert.Empty(t, r.URL.Fragment)
		return &http.Response{
			StatusCode: status,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
			Request:    r,
		}, nil
	})
	return &secretDataSourceTestConfig{client: &gophercloud.ServiceClient{
		ProviderClient: provider,
		Endpoint:       "https://kms.example.test/kms/user/v1/",
	}}
}

func TestAccKMSKeyDataSource_basic(t *testing.T) {
	name := "tfacc-kms-key-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccKMSPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{{Config: fmt.Sprintf(testAccKMSKeyDataSourceBasic, name, "aes256-gcm96", true), Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("data.vkcs_kms_key.key", "id", name),
			resource.TestCheckResourceAttrPair("data.vkcs_kms_key.key", "type", "vkcs_kms_key.key", "type"),
			resource.TestCheckResourceAttrPair("data.vkcs_kms_key.key", "deletion_allowed", "vkcs_kms_key.key", "deletion_allowed"),
		)}},
	})
}

const testAccKMSKeyDataSourceBasic = testAccKMSKeyResourceConfig + `
data "vkcs_kms_key" "key" {
  name = vkcs_kms_key.key.name
}
`
