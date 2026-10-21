package kms_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"
)

func TestKMSKeyDecryptDataSourceRead(t *testing.T) {
	config := newKeyDecryptDataSourceTestConfig(t, http.StatusOK, `{"data":{"plaintext":"aGVsbG8="}}`, nil)
	ds := kms.NewKeyDecryptDataSource()
	var schema datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{
		kms.KeyCryptFieldKey: "test-key", kms.KeyCryptFieldCiphertext: "vault:v1:dGVzdA==", kms.KeyCryptFieldPlaintext: "stale-output",
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
	assert.Equal(t, "333d6b3a3c1f5db6c9bdda5939b136986d170f4649172a68368d54ecb44c2ff2", attribute1.ValueString())
	var attribute2 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeyCryptFieldPlaintext), &attribute2).HasError())
	assert.Equal(t, "aGVsbG8=", attribute2.ValueString())
	var attribute3 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeyCryptFieldKey), &attribute3).HasError())
	assert.Equal(t, "test-key", attribute3.ValueString())
	var attribute4 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeyCryptFieldCiphertext), &attribute4).HasError())
	assert.Equal(t, "vault:v1:dGVzdA==", attribute4.ValueString())
}

func TestKMSKeyDecryptDataSourceRead_errors(t *testing.T) {
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
		{name: "invalid output type", status: http.StatusOK, body: `{"data":{"plaintext":123}}`, wantError: "cannot unmarshal"},
		{name: "transport error", transportErr: errors.New("connection failed"), wantError: "connection failed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newKeyDecryptDataSourceTestConfig(t, tc.status, tc.body, tc.transportErr)
			ds := kms.NewKeyDecryptDataSource()
			var schema datasource.SchemaResponse
			ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{
				kms.KeyCryptFieldKey: "test-key", kms.KeyCryptFieldCiphertext: "vault:v1:dGVzdA==",
			})

			ds.(datasource.DataSourceWithConfigure).Configure(context.Background(), datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
			response := datasource.ReadResponse{State: d}
			ds.Read(context.Background(), datasource.ReadRequest{Config: tfsdk.Config{Schema: d.Schema, Raw: d.Raw}}, &response)
			d = response.State
			diags := response.Diagnostics
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary(), "Error decrypting")
			assert.Contains(t, diags[0].Detail(), tc.wantError)
			var attribute5 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute5).HasError())
			assert.Empty(t, attribute5.ValueString())
			var attribute6 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeyCryptFieldPlaintext), &attribute6).HasError())
			assert.Empty(t, attribute6.ValueString())
		})
	}
}

func TestKMSKeyDecryptDataSourceRead_clientError(t *testing.T) {
	ds := kms.NewKeyDecryptDataSource()
	var schema datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{
		kms.KeyCryptFieldKey: "test-key", kms.KeyCryptFieldCiphertext: "vault:v1:dGVzdA==",
	})
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
	var attribute7 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute7).HasError())
	assert.Empty(t, attribute7.ValueString())
	var attribute8 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeyCryptFieldPlaintext), &attribute8).HasError())
	assert.Empty(t, attribute8.ValueString())
}

func newKeyDecryptDataSourceTestConfig(t *testing.T, status int, body string, transportErr error) *secretDataSourceTestConfig {
	t.Helper()
	calls := 0
	t.Cleanup(func() { assert.Equal(t, 1, calls, "expected one KMS request") })
	provider := &gophercloud.ProviderClient{}
	provider.HTTPClient.Transport = secretDataSourceTestTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/kms/user/v1/transit/decrypt/test-key", r.URL.Path)
		requestBody, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `{"ciphertext":"vault:v1:dGVzdA=="}`, string(requestBody))
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
