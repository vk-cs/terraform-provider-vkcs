package kms_test

import (
	"context"
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/require"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"
)

func TestKMSReadRetries(t *testing.T) {
	for _, tc := range []struct {
		name     string
		endpoint any
		input    map[string]interface{}
		requests []kmsResourceRequest
	}{
		{"key resource", kms.NewKeyResource(), map[string]interface{}{kms.KeyFieldName: "test-key"}, []kmsResourceRequest{{"GET", "transit/keys/test-key", "", 200, `{"data":{"name":"test-key","type":"aes256-gcm96"}}`}}},
		{"key data source", kms.NewKeyDataSource(), map[string]interface{}{kms.KeyFieldName: "test-key"}, []kmsResourceRequest{{"GET", "transit/keys/test-key", "", 200, `{"data":{"name":"test-key","type":"aes256-gcm96"}}`}}},
		{"secret resource", kms.NewSecretResource(), map[string]interface{}{kms.SecretFieldPath: "test-secret"}, secretReadRetryRequests()},
		{"secret data source", kms.NewSecretDataSource(), map[string]interface{}{kms.SecretFieldPath: "test-secret"}, secretReadRetryRequests()},
		{"keys list", kms.NewKeysListDataSource(), nil, []kmsResourceRequest{{"GET", "transit/keys", "", 200, `{"data":{"keys":["test-key"]}}`}}},
		{"secrets list", kms.NewSecretsListDataSource(), nil, []kmsResourceRequest{{"GET", "secret/metadata/", "", 200, `{"data":{"keys":["test-secret"]}}`}}},
		{"encrypt", kms.NewKeyEncryptDataSource(), map[string]interface{}{kms.KeyCryptFieldKey: "test-key", kms.KeyCryptFieldPlaintext: "aGVsbG8="}, []kmsResourceRequest{{"POST", "transit/encrypt/test-key", `{"plaintext":"aGVsbG8="}`, 200, `{"data":{"ciphertext":"vault:v1:test"}}`}}},
		{"decrypt", kms.NewKeyDecryptDataSource(), map[string]interface{}{kms.KeyCryptFieldKey: "test-key", kms.KeyCryptFieldCiphertext: "vault:v1:test"}, []kmsResourceRequest{{"POST", "transit/decrypt/test-key", `{"ciphertext":"vault:v1:test"}`, 200, `{"data":{"plaintext":"aGVsbG8="}}`}}},
	} {
		for _, status := range []int{400, 404} {
			t.Run(tc.name+map[int]string{400: "/unfinished", 404: "/not found"}[status], func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					var requests []kmsResourceRequest
					for _, success := range tc.requests {
						retry := success
						retry.status = status
						retry.response = `{"errors":["unfinished operation"]}`
						requests = append(requests, retry, success)
					}
					endpoint := newFrameworkTestEndpoint(t, tc.endpoint)
					data := newFrameworkTestData(t, endpoint, tc.input)
					data.SetId("existing")
					diags := endpoint.ReadContext(context.Background(), data, newKMSResourceTestConfig(t, requests...))
					require.False(t, diags.HasError(), "%v", diags)
				})
			})
		}
	}
}

func secretReadRetryRequests() []kmsResourceRequest {
	return []kmsResourceRequest{
		{"GET", "secret/data/test-secret", "", 200, `{"data":{"data":{"value":"test"},"metadata":{"version":1}}}`},
		{"GET", "secret/metadata/delete-protection/test-secret", "", 200, `{"delete_protection":true}`},
	}
}
