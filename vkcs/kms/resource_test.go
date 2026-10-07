package kms_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"

	"github.com/gophercloud/gophercloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		assert.Empty(t, r.URL.RawQuery)
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
	for name, r := range map[string]*schema.Resource{"key": kms.ResourceKey(), "secret": kms.ResourceSecret()} {
		for operation, call := range map[string]func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics{"create": r.CreateContext, "update": r.UpdateContext, "delete": r.DeleteContext} {
			t.Run(name+"/"+operation, func(t *testing.T) {
				d := schema.TestResourceDataRaw(t, r.Schema, nil)
				if operation != "create" {
					d.SetId("existing")
				}
				id := d.Id()
				diags := call(context.Background(), d, &secretDataSourceTestConfig{err: errors.New("client unavailable")})
				require.True(t, diags.HasError())
				require.Len(t, diags, 1)
				assert.Equal(t, "Error creating VKCS KMS client: client unavailable", diags[0].Summary)
				assert.Equal(t, id, d.Id())
			})
		}
	}
}
