package kms_test

import (
	"context"
	"testing"
	"testing/synctest"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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
		{"key resource", kms.NewKeyResource(), map[string]interface{}{"name": "test-key"}, []kmsResourceRequest{{"GET", "transit/keys/test-key", "", 200, `{"data":{"name":"test-key","type":"aes256-gcm96"}}`}}},
		{"key data source", kms.NewKeyDataSource(), map[string]interface{}{"name": "test-key"}, []kmsResourceRequest{{"GET", "transit/keys/test-key", "", 200, `{"data":{"name":"test-key","type":"aes256-gcm96"}}`}}},
		{"secret resource", kms.NewSecretResource(), map[string]interface{}{"path": "test-secret"}, secretReadRetryRequests()},
		{"secret data source", kms.NewSecretDataSource(), map[string]interface{}{"path": "test-secret"}, secretReadRetryRequests()},
		{"keys list", kms.NewKeysListDataSource(), nil, []kmsResourceRequest{{"GET", "transit/keys", "", 200, `{"data":{"keys":["test-key"]}}`}}},
		{"secrets list", kms.NewSecretsListDataSource(), nil, []kmsResourceRequest{{"GET", "secret/metadata/", "", 200, `{"data":{"keys":["test-secret"]}}`}}},
		{"encrypt", kms.NewKeyEncryptDataSource(), map[string]interface{}{"key": "test-key", "plaintext": "aGVsbG8="}, []kmsResourceRequest{{"POST", "transit/encrypt/test-key", `{"plaintext":"aGVsbG8="}`, 200, `{"data":{"ciphertext":"vault:v1:test"}}`}}},
		{"decrypt", kms.NewKeyDecryptDataSource(), map[string]interface{}{"key": "test-key", "ciphertext": "vault:v1:test"}, []kmsResourceRequest{{"POST", "transit/decrypt/test-key", `{"ciphertext":"vault:v1:test"}`, 200, `{"data":{"plaintext":"aGVsbG8="}}`}}},
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
					ctx := context.Background()
					config := newKMSResourceTestConfig(t, requests...)
					var diags diag.Diagnostics
					switch endpoint := tc.endpoint.(type) {
					case resource.ResourceWithConfigure:
						var schema resource.SchemaResponse
						endpoint.Schema(ctx, resource.SchemaRequest{}, &schema)
						state := newTestState(t, tfsdk.State{Schema: schema.Schema}, tc.input)
						require.False(t, state.SetAttribute(ctx, path.Root("id"), "existing").HasError())
						endpoint.Configure(ctx, resource.ConfigureRequest{ProviderData: config}, &resource.ConfigureResponse{})
						response := resource.ReadResponse{State: state}
						endpoint.Read(ctx, resource.ReadRequest{State: state}, &response)
						diags = response.Diagnostics
					case datasource.DataSourceWithConfigure:
						var schema datasource.SchemaResponse
						endpoint.Schema(ctx, datasource.SchemaRequest{}, &schema)
						state := newTestState(t, tfsdk.State{Schema: schema.Schema}, tc.input)
						endpoint.Configure(ctx, datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
						response := datasource.ReadResponse{State: state}
						endpoint.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: state.Schema, Raw: state.Raw}}, &response)
						diags = response.Diagnostics
					}
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
