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

func TestKMSKeyEncryptDataSourceRead(t *testing.T) {
	config := newKeyEncryptDataSourceTestConfig(t, http.StatusOK, `{"data":{"ciphertext":"vault:v1:dGVzdA=="}}`, nil)
	ds := kms.NewKeyEncryptDataSource()
	var schema datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{
		"key": "test-key", "plaintext": "aGVsbG8=", "ciphertext": "stale-output",
	})
	require.False(t, d.SetAttribute(context.Background(), path.Root("id"), "stale-id").HasError())

	ds.(datasource.DataSourceWithConfigure).Configure(context.Background(), datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
	response := datasource.ReadResponse{State: d}
	ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config(d)}, &response)
	d = response.State
	diags := response.Diagnostics
	require.False(t, diags.HasError(), "%v", diags)
	var attribute1 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root("id"), &attribute1).HasError())
	assert.Equal(t, "655d0adeea5986d6235be9f53d707ae305422728b3038d1ab77d397a8cfe69d2", attribute1.ValueString())
	var attribute2 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root("ciphertext"), &attribute2).HasError())
	assert.Equal(t, "vault:v1:dGVzdA==", attribute2.ValueString())
	var attribute3 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root("key"), &attribute3).HasError())
	assert.Equal(t, "test-key", attribute3.ValueString())
	var attribute4 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root("plaintext"), &attribute4).HasError())
	assert.Equal(t, "aGVsbG8=", attribute4.ValueString())
}

func TestKMSKeyEncryptDataSourceRead_errors(t *testing.T) {
	for _, tc := range []struct {
		name         string
		status       int
		body         string
		transportErr error
		wantError    string
	}{
		{name: "forbidden", status: http.StatusForbidden, body: `{"errors":["permission denied"]}`, wantError: "permission denied"},
		{name: "bad request", status: http.StatusBadRequest, body: `{"errors":["invalid input"]}`, wantError: "invalid input"},
		{name: "invalid JSON", status: http.StatusOK, body: `invalid`, wantError: "invalid character"},
		{name: "invalid output type", status: http.StatusOK, body: `{"data":{"ciphertext":123}}`, wantError: "cannot unmarshal"},
		{name: "transport error", transportErr: errors.New("connection failed"), wantError: "connection failed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newKeyEncryptDataSourceTestConfig(t, tc.status, tc.body, tc.transportErr)
			ds := kms.NewKeyEncryptDataSource()
			var schema datasource.SchemaResponse
			ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{
				"key": "test-key", "plaintext": "aGVsbG8=",
			})

			ds.(datasource.DataSourceWithConfigure).Configure(context.Background(), datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
			response := datasource.ReadResponse{State: d}
			ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config(d)}, &response)
			d = response.State
			diags := response.Diagnostics
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary(), "Error encrypting")
			assert.Contains(t, diags[0].Detail(), tc.wantError)
			var attribute5 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root("id"), &attribute5).HasError())
			assert.Empty(t, attribute5.ValueString())
			var attribute6 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root("ciphertext"), &attribute6).HasError())
			assert.Empty(t, attribute6.ValueString())
		})
	}
}

func TestKMSKeyEncryptDataSourceRead_clientError(t *testing.T) {
	ds := kms.NewKeyEncryptDataSource()
	var schema datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{
		"key": "test-key", "plaintext": "aGVsbG8=",
	})
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
	var attribute7 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root("id"), &attribute7).HasError())
	assert.Empty(t, attribute7.ValueString())
	var attribute8 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root("ciphertext"), &attribute8).HasError())
	assert.Empty(t, attribute8.ValueString())
}

func newKeyEncryptDataSourceTestConfig(t *testing.T, status int, body string, transportErr error) *secretDataSourceTestConfig {
	t.Helper()
	calls := 0
	t.Cleanup(func() { assert.Equal(t, 1, calls, "expected one KMS request") })
	provider := &gophercloud.ProviderClient{}
	provider.HTTPClient.Transport = secretDataSourceTestTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/kms/user/v1/transit/encrypt/test-key", r.URL.Path)
		requestBody, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `{"plaintext":"aGVsbG8="}`, string(requestBody))
		if transportErr != nil {
			return nil, transportErr
		}
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

// Both cryptographic data sources are exercised together to verify a real round trip.
func TestAccKMSKeyEncryptDecryptDataSources_basic(t *testing.T) {
	name := "tfacc-kms-key-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	config := fmt.Sprintf(testAccKMSKeyEncryptDecryptDataSourcesBasic, name, "aes256-gcm96", true)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{{Config: config, Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet("data.vkcs_kms_key_encrypt.encrypted", "id"),
			resource.TestCheckResourceAttrSet("data.vkcs_kms_key_encrypt.encrypted", "ciphertext"),
			resource.TestCheckResourceAttrSet("data.vkcs_kms_key_decrypt.decrypted", "id"),
			resource.TestCheckResourceAttrPair("data.vkcs_kms_key_decrypt.decrypted", "plaintext", "data.vkcs_kms_key_encrypt.encrypted", "plaintext"),
		)}},
	})
}

const testAccKMSKeyEncryptDecryptDataSourcesBasic = testAccKMSKeyResourceConfig + `
data "vkcs_kms_key_encrypt" "encrypted" {
  key       = vkcs_kms_key.key.name
  plaintext = base64encode("acceptance-test plaintext")
}

data "vkcs_kms_key_decrypt" "decrypted" {
  key        = vkcs_kms_key.key.name
  ciphertext = data.vkcs_kms_key_encrypt.encrypted.ciphertext
}
`
