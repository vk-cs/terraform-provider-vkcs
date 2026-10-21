package kms_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"
)

func TestKMSSecretDataSourceRead(t *testing.T) {
	body := `{"data":{"data":{"username":"alice","password":"test-password"},"metadata":{"created_time":"2026-09-01T10:00:00Z","version":3}}}`
	config := newSecretDataSourceTestConfig(t, http.StatusOK, body, true)
	ds := kms.NewSecretDataSource()
	var schema datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.SecretFieldPath: "test-secret", kms.SecretFieldDataJSON: `{"stale":"value"}`, kms.SecretFieldCreatedTime: "stale-time", kms.SecretFieldVersion: 1})

	require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "stale-id").HasError())
	ds.(datasource.DataSourceWithConfigure).Configure(context.Background(), datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
	response := datasource.ReadResponse{State: d}
	ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config{Schema: d.Schema, Raw: d.Raw}}, &response)
	d = response.State
	diags := response.Diagnostics
	require.False(t, diags.HasError(), "%v", diags)
	var attribute1 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute1).HasError())
	assert.Equal(t, "test-secret", attribute1.ValueString())
	var attribute2 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldPath), &attribute2).HasError())
	assert.Equal(t, "test-secret", attribute2.ValueString())
	var attribute3 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldDataJSON), &attribute3).HasError())
	assert.JSONEq(t, `{"username":"alice","password":"test-password"}`, attribute3.ValueString())
	var attribute4 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldCreatedTime), &attribute4).HasError())
	assert.Equal(t, "2026-09-01 10:00:00 +0000 UTC", attribute4.ValueString())
	var attribute5 types.Int64
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldVersion), &attribute5).HasError())
	assert.Equal(t, 3, int(attribute5.ValueInt64()))
	var attribute6 types.Bool
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldDeleteProtection), &attribute6).HasError())
	assert.Equal(t, true, attribute6.ValueBool())
}

func TestKMSSecretDataSourceRead_errors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{name: "forbidden", status: http.StatusForbidden, body: `{"errors":["permission denied"]}`},
		{name: "invalid JSON", status: http.StatusOK, body: `invalid`},
		{name: "invalid data type", status: http.StatusOK, body: `{"data":{"data":{"count":123}}}`},
		{name: "invalid timestamp", status: http.StatusOK, body: `{"data":{"metadata":{"created_time":"invalid"}}}`},
		{name: "invalid version", status: http.StatusOK, body: `{"data":{"metadata":{"version":"invalid"}}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newSecretDataSourceTestConfig(t, tc.status, tc.body)
			ds := kms.NewSecretDataSource()
			var schema datasource.SchemaResponse
			ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.SecretFieldPath: "test-secret"})

			ds.(datasource.DataSourceWithConfigure).Configure(context.Background(), datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
			response := datasource.ReadResponse{State: d}
			ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config{Schema: d.Schema, Raw: d.Raw}}, &response)
			d = response.State
			diags := response.Diagnostics
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary(), "Error getting secret")
			var attribute7 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute7).HasError())
			assert.Empty(t, attribute7.ValueString())
		})
	}
}

func TestKMSSecretDataSourceRead_clientError(t *testing.T) {
	ds := kms.NewSecretDataSource()
	var schema datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.SecretFieldPath: "test-secret"})
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
	var attribute8 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute8).HasError())
	assert.Empty(t, attribute8.ValueString())
}

type secretDataSourceTestConfig struct {
	clients.Config
	client *gophercloud.ServiceClient
	err    error
}

func (c *secretDataSourceTestConfig) GetRegion() string { return "RegionOne" }

func (c *secretDataSourceTestConfig) KMSV1Client(string) (*gophercloud.ServiceClient, error) {
	return c.client, c.err
}

type secretDataSourceTestTransport func(*http.Request) (*http.Response, error)

func (f secretDataSourceTestTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func newSecretDataSourceTestConfig(t *testing.T, status int, body string, protection ...bool) *secretDataSourceTestConfig {
	t.Helper()
	requests := []kmsResourceRequest{{"GET", "secret/data/test-secret", "", status, body}}
	for _, enabled := range protection {
		requests = append(requests, kmsResourceRequest{"GET", "secret/metadata/delete-protection/test-secret", "", 200, fmt.Sprintf(`{"delete_protection":%t}`, enabled)})
	}
	return newKMSResourceTestConfig(t, requests...)
}

func TestKMSSecretDataSourceRead_escapedPath(t *testing.T) {
	ds := kms.NewSecretDataSource()
	var schema datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.SecretFieldPath: "secret /?#%"})
	config := newKMSResourceTestConfig(t, kmsResourceRequest{
		method:   http.MethodGet,
		path:     "secret/data/secret%20%2F%3F%23%25",
		status:   http.StatusOK,
		response: `{"data":{"data":{"password":"test-password"},"metadata":{"created_time":"2026-09-01T10:00:00Z","version":1}}}`,
	}, kmsResourceRequest{"GET", "secret/metadata/delete-protection/secret%20%2F%3F%23%25", "", 200, `{"delete_protection":false}`})
	ds.(datasource.DataSourceWithConfigure).Configure(context.Background(), datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
	response := datasource.ReadResponse{State: d}
	ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config{Schema: d.Schema, Raw: d.Raw}}, &response)
	d = response.State
	diags := response.Diagnostics
	require.False(t, diags.HasError(), "%v", diags)
	var attribute9 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute9).HasError())
	assert.Equal(t, "secret /?#%", attribute9.ValueString())
	var attribute10 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldPath), &attribute10).HasError())
	assert.Equal(t, "secret /?#%", attribute10.ValueString())
	var attribute11 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldDataJSON), &attribute11).HasError())
	assert.JSONEq(t, `{"password":"test-password"}`, attribute11.ValueString())
}

func TestAccKMSSecretDataSource_basic(t *testing.T) {
	name := "tfacc-kms-secret-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccKMSPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKMSSecretDataSourceBasic, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_kms_secret.secret", kms.EntityFieldID, name),
					resource.TestCheckResourceAttr("data.vkcs_kms_secret.secret", kms.SecretFieldPath, name),
					resource.TestCheckResourceAttrPair("data.vkcs_kms_secret.secret", kms.SecretFieldDataJSON, "vkcs_kms_secret.secret", kms.SecretFieldDataJSON),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secret.secret", kms.SecretFieldCreatedTime),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secret.secret", kms.SecretFieldVersion),
				),
			},
		},
	})
}

const testAccKMSSecretDataSourceBasic = `
resource "vkcs_kms_secret" "secret" {
  path              = %q
  delete_protection = false
  data_json         = jsonencode({ password = "test-password" })
}

data "vkcs_kms_secret" "secret" {
  path = vkcs_kms_secret.secret.path
}
`
