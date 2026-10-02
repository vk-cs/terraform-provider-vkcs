package kms_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"
)

func TestKMSKeyDecryptDataSourceRead(t *testing.T) {
	config := newKeyDecryptDataSourceTestConfig(t, http.StatusOK, `{"data":{"plaintext":"aGVsbG8="}}`, nil)
	ds := kms.DataSourceKeyDecrypt()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"key": "test-key", "ciphertext": "vault:v1:dGVzdA==", "plaintext": "stale-output",
	})
	d.SetId("stale-id")

	diags := ds.ReadContext(context.Background(), d, config)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "333d6b3a3c1f5db6c9bdda5939b136986d170f4649172a68368d54ecb44c2ff2", d.Id())
	assert.Equal(t, "aGVsbG8=", d.Get("plaintext"))
	assert.Equal(t, "test-key", d.Get("key"))
	assert.Equal(t, "vault:v1:dGVzdA==", d.Get("ciphertext"))
}

func TestKMSKeyDecryptDataSourceRead_errors(t *testing.T) {
	for _, tc := range []struct {
		name         string
		status       int
		body         string
		transportErr error
		wantError    string
	}{
		{name: "not found", status: http.StatusNotFound, body: `{"errors":["key not found"]}`, wantError: "key not found"},
		{name: "forbidden", status: http.StatusForbidden, body: `{"errors":["permission denied"]}`, wantError: "permission denied"},
		{name: "bad request", status: http.StatusBadRequest, body: `{"errors":["invalid input"]}`, wantError: "invalid input"},
		{name: "invalid JSON", status: http.StatusOK, body: `invalid`, wantError: "invalid character"},
		{name: "invalid output type", status: http.StatusOK, body: `{"data":{"plaintext":123}}`, wantError: "cannot unmarshal"},
		{name: "transport error", transportErr: errors.New("connection failed"), wantError: "connection failed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newKeyDecryptDataSourceTestConfig(t, tc.status, tc.body, tc.transportErr)
			ds := kms.DataSourceKeyDecrypt()
			d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
				"key": "test-key", "ciphertext": "vault:v1:dGVzdA==",
			})

			diags := ds.ReadContext(context.Background(), d, config)
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary, "Error decrypting:")
			assert.Contains(t, diags[0].Summary, tc.wantError)
			assert.Empty(t, d.Id())
			assert.Empty(t, d.Get("plaintext"))
		})
	}
}

func TestKMSKeyDecryptDataSourceRead_clientError(t *testing.T) {
	ds := kms.DataSourceKeyDecrypt()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"key": "test-key", "ciphertext": "vault:v1:dGVzdA==",
	})
	config := &secretDataSourceTestConfig{err: errors.New("client unavailable")}

	diags := ds.ReadContext(context.Background(), d, config)
	require.True(t, diags.HasError())
	require.Len(t, diags, 1)
	assert.Equal(t, "Error creating VKCS KMS client: client unavailable", diags[0].Summary)
	assert.Empty(t, d.Id())
	assert.Empty(t, d.Get("plaintext"))
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
