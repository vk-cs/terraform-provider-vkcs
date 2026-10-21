package kms_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"
)

type kmsResourceRequest struct {
	method, path, body string
	status             int
	response           string
}

func newKMSResourceTestConfig(t *testing.T, requests ...kmsResourceRequest) *secretDataSourceTestConfig {
	t.Helper()
	next := 0
	t.Cleanup(func() { assert.Equal(t, len(requests), next, "request count") })
	provider := &gophercloud.ProviderClient{}
	provider.HTTPClient.Transport = secretDataSourceTestTransport(func(r *http.Request) (*http.Response, error) {
		require.Less(t, next, len(requests), "unexpected request: %s %s", r.Method, r.URL)
		expected := requests[next]
		next++
		assert.Equal(t, expected.method, r.Method)
		assert.Equal(t, "/kms/user/v1/"+expected.path, r.URL.EscapedPath())
		if expected.path == "transit/keys" || expected.path == "secret/metadata/" {
			assert.Equal(t, "list=true", r.URL.RawQuery)
		} else {
			assert.Empty(t, r.URL.RawQuery)
		}
		var body []byte
		if r.Body != nil {
			var err error
			body, err = io.ReadAll(r.Body)
			require.NoError(t, err)
		}
		if expected.body == "" {
			assert.Empty(t, body)
		} else {
			assert.JSONEq(t, expected.body, string(body))
		}
		return &http.Response{StatusCode: expected.status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(expected.response)), Request: r}, nil
	})
	return &secretDataSourceTestConfig{client: &gophercloud.ServiceClient{ProviderClient: provider, Endpoint: "https://kms.example.test/kms/user/v1/"}}
}

func TestKMSResourceWriteClientErrors(t *testing.T) {
	for name, r := range map[string]resource.Resource{"key": kms.NewKeyResource(), "secret": kms.NewSecretResource()} {
		for _, operation := range []string{"create", "update", "delete"} {
			t.Run(name+"/"+operation, func(t *testing.T) {
				ctx := context.Background()
				var schema resource.SchemaResponse
				r.Schema(ctx, resource.SchemaRequest{}, &schema)
				state := newTestState(t, tfsdk.State{Schema: schema.Schema}, nil)
				if operation != "create" {
					require.False(t, state.SetAttribute(ctx, path.Root(kms.EntityFieldID), "existing").HasError())
				}
				r.(resource.ResourceWithConfigure).Configure(ctx, resource.ConfigureRequest{ProviderData: &secretDataSourceTestConfig{err: errors.New("client unavailable")}}, &resource.ConfigureResponse{})
				originalState := state.Raw
				var diags diag.Diagnostics
				switch operation {
				case "create":
					response := resource.CreateResponse{State: state}
					r.Create(ctx, resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}}, &response)
					diags = response.Diagnostics
					state = response.State
				case "update":
					response := resource.UpdateResponse{State: state}
					r.Update(ctx, resource.UpdateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}}, &response)
					diags = response.Diagnostics
					state = response.State
				case "delete":
					response := resource.DeleteResponse{State: state}
					r.Delete(ctx, resource.DeleteRequest{State: state}, &response)
					diags = response.Diagnostics
					state = response.State
				}
				require.True(t, diags.HasError())
				require.Len(t, diags, 1)
				assert.Equal(t, "Error creating VKCS KMS client", diags[0].Summary())
				assert.Equal(t, "client unavailable", diags[0].Detail())
				assert.True(t, originalState.Equal(state.Raw), "state must remain unchanged on client errors")
			})
		}
	}
}
