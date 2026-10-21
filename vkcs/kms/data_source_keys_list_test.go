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

func TestKMSKeysListDataSourceRead(t *testing.T) {
	for _, tc := range []struct {
		name     string
		body     string
		wantKeys []interface{}
		wantID   string
	}{
		{
			name:     "multiple keys",
			body:     `{"data":{"keys":["test-key","another-key"]}}`,
			wantKeys: []interface{}{"test-key", "another-key"},
			wantID:   "eb2e78d667d7f3ed30a8c3a93174895e2fc91e317acdbb28d7f53250114bfc07",
		},
		{
			name:     "empty list",
			body:     `{"data":{"keys":[]}}`,
			wantKeys: []interface{}{},
			wantID:   "19e89348f2a9d5f3d0c5fca8e2a7068d9c5d71687a355f759009a5ed2527eb2c",
		},
		{
			name:     "null list",
			body:     `{"data":{"keys":null}}`,
			wantKeys: []interface{}{},
			wantID:   "19e89348f2a9d5f3d0c5fca8e2a7068d9c5d71687a355f759009a5ed2527eb2c",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newKeysListDataSourceTestConfig(t, http.StatusOK, tc.body)
			ds := kms.NewKeysListDataSource()
			var schema datasource.SchemaResponse
			ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{
				kms.KeysListFieldKeys: []interface{}{"stale-key"}, kms.KeysListFieldKeysCount: 1,
			})

			ds.(datasource.DataSourceWithConfigure).Configure(context.Background(), datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
			response := datasource.ReadResponse{State: d}
			ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config{Schema: d.Schema, Raw: d.Raw}}, &response)
			d = response.State
			diags := response.Diagnostics
			require.False(t, diags.HasError(), "%v", diags)
			var attribute1 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute1).HasError())
			assert.Equal(t, tc.wantID, attribute1.ValueString())
			var attribute2 types.Int64
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeysListFieldKeysCount), &attribute2).HasError())
			assert.Equal(t, len(tc.wantKeys), int(attribute2.ValueInt64()))
			var attribute3 []string
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeysListFieldKeys), &attribute3).HasError())
			assert.ElementsMatch(t, tc.wantKeys, attribute3)
		})
	}
}

func TestKMSKeysListDataSourceRead_errors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{name: "forbidden", status: http.StatusForbidden, body: `{"errors":["permission denied"]}`},
		{name: "invalid JSON", status: http.StatusOK, body: `invalid`},
		{name: "invalid keys type", status: http.StatusOK, body: `{"data":{"keys":123}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newKeysListDataSourceTestConfig(t, tc.status, tc.body)
			ds := kms.NewKeysListDataSource()
			var schema datasource.SchemaResponse
			ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, nil)

			ds.(datasource.DataSourceWithConfigure).Configure(context.Background(), datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
			response := datasource.ReadResponse{State: d}
			ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config{Schema: d.Schema, Raw: d.Raw}}, &response)
			d = response.State
			diags := response.Diagnostics
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary(), "Error listing keys")
			var attribute4 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute4).HasError())
			assert.Empty(t, attribute4.ValueString())
		})
	}
}

func TestKMSKeysListDataSourceRead_clientError(t *testing.T) {
	ds := kms.NewKeysListDataSource()
	var schema datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, nil)
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
	var attribute5 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute5).HasError())
	assert.Empty(t, attribute5.ValueString())
}

func newKeysListDataSourceTestConfig(t *testing.T, status int, body string) *secretDataSourceTestConfig {
	t.Helper()
	calls := 0
	t.Cleanup(func() { assert.Equal(t, 1, calls) })
	provider := &gophercloud.ProviderClient{}
	provider.HTTPClient.Transport = secretDataSourceTestTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/kms/user/v1/transit/keys", r.URL.Path)
		assert.Equal(t, "list=true", r.URL.RawQuery)
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

func TestAccKMSKeysListDataSource_basic(t *testing.T) {
	name := "tfacc-kms-key-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccKMSPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{{Config: fmt.Sprintf(testAccKMSKeysListDataSourceBasic, name, "aes256-gcm96", true), Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet("data.vkcs_kms_keys_list.keys", "id"),
			resource.TestCheckResourceAttrSet("data.vkcs_kms_keys_list.keys", "keys_count"),
			resource.TestCheckTypeSetElemAttr("data.vkcs_kms_keys_list.keys", "keys.*", name),
		)}},
	})
}

const testAccKMSKeysListDataSourceBasic = testAccKMSKeyResourceConfig + `
data "vkcs_kms_keys_list" "keys" {
  depends_on = [vkcs_kms_key.key]
}
`
