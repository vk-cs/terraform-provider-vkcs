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

func TestKMSSecretsDataSourceRead(t *testing.T) {
	for _, tc := range []struct {
		name        string
		body        string
		wantSecrets []interface{}
		wantID      string
	}{
		{
			name:        "multiple secrets",
			body:        `{"data":{"keys":["test-secret","another-secret"]}}`,
			wantSecrets: []interface{}{"test-secret", "another-secret"},
			wantID:      "7a9fbe1cce12c710a5bdc4a5e3bd05de40d59ddf96b6104ea0376bc5f6fe485f",
		},
		{
			name:        "empty list",
			body:        `{"data":{"keys":[]}}`,
			wantSecrets: []interface{}{},
			wantID:      "19e89348f2a9d5f3d0c5fca8e2a7068d9c5d71687a355f759009a5ed2527eb2c",
		},
		{
			name:        "null list",
			body:        `{"data":{"keys":null}}`,
			wantSecrets: []interface{}{},
			wantID:      "19e89348f2a9d5f3d0c5fca8e2a7068d9c5d71687a355f759009a5ed2527eb2c",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newSecretsDataSourceTestConfig(t, http.StatusOK, tc.body)
			ds := kms.NewSecretsDataSource()
			var schema datasource.SchemaResponse
			ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{
				"secrets": []interface{}{"stale-secret"}, "secrets_count": 1,
			})

			ds.(datasource.DataSourceWithConfigure).Configure(context.Background(), datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
			response := datasource.ReadResponse{State: d}
			ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config(d)}, &response)
			d = response.State
			diags := response.Diagnostics
			require.False(t, diags.HasError(), "%v", diags)
			var attribute1 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root("id"), &attribute1).HasError())
			assert.Equal(t, tc.wantID, attribute1.ValueString())
			var attribute2 types.Int64
			require.False(t, d.GetAttribute(context.Background(), path.Root("secrets_count"), &attribute2).HasError())
			assert.Equal(t, len(tc.wantSecrets), int(attribute2.ValueInt64()))
			var attribute3 []string
			require.False(t, d.GetAttribute(context.Background(), path.Root("secrets"), &attribute3).HasError())
			assert.ElementsMatch(t, tc.wantSecrets, attribute3)
		})
	}
}

func TestKMSSecretsDataSourceRead_errors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{name: "forbidden", status: http.StatusForbidden, body: `{"errors":["permission denied"]}`},
		{name: "invalid JSON", status: http.StatusOK, body: `invalid`},
		{name: "invalid keys type", status: http.StatusOK, body: `{"data":{"keys":123}}`},
		{name: "invalid key element", status: http.StatusOK, body: `{"data":{"keys":[123]}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newSecretsDataSourceTestConfig(t, tc.status, tc.body)
			ds := kms.NewSecretsDataSource()
			var schema datasource.SchemaResponse
			ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, nil)

			ds.(datasource.DataSourceWithConfigure).Configure(context.Background(), datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
			response := datasource.ReadResponse{State: d}
			ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config(d)}, &response)
			d = response.State
			diags := response.Diagnostics
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary(), "Error listing secrets")
			var attribute4 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root("id"), &attribute4).HasError())
			assert.Empty(t, attribute4.ValueString())
		})
	}
}

func TestKMSSecretsDataSourceRead_clientError(t *testing.T) {
	ds := kms.NewSecretsDataSource()
	var schema datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, nil)
	config := &secretDataSourceTestConfig{err: errors.New("client unavailable")}

	ds.(datasource.DataSourceWithConfigure).Configure(context.Background(), datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
	response := datasource.ReadResponse{State: d}
	ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config(d)}, &response)
	d = response.State
	diags := response.Diagnostics
	require.True(t, diags.HasError())
	require.Len(t, diags, 1)
	assert.Equal(t, "Error creating VKCS KMS client", diags[0].Summary())
	assert.Equal(t, "client unavailable", diags[0].Detail())
	var attribute5 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root("id"), &attribute5).HasError())
	assert.Empty(t, attribute5.ValueString())
}

func newSecretsDataSourceTestConfig(t *testing.T, status int, body string) *secretDataSourceTestConfig {
	t.Helper()
	calls := 0
	t.Cleanup(func() { assert.Equal(t, 1, calls, "expected one KMS request") })
	provider := &gophercloud.ProviderClient{}
	provider.HTTPClient.Transport = secretDataSourceTestTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		assert.Empty(t, r.URL.Fragment)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/kms/user/v1/secret/metadata/", r.URL.Path)
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

func TestAccKMSSecretsDataSource_basic(t *testing.T) {
	name := "tfacc-kms-secret-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKMSSecretsDataSourceBasic, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secrets.secrets", "id"),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secrets.secrets", "secrets_count"),
					resource.TestCheckTypeSetElemAttr("data.vkcs_kms_secrets.secrets", "secrets.*", name),
				),
			},
		},
	})
}

const testAccKMSSecretsDataSourceBasic = `
resource "vkcs_kms_secret" "secret" {
  path              = %q
  delete_protection = false
  data_json         = jsonencode({ password = "test-password" })
}

data "vkcs_kms_secrets" "secrets" {
  depends_on = [vkcs_kms_secret.secret]
}
`
