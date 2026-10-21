package kms

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/synctest"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/stretchr/testify/require"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/keys"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/secrets"
)

type refreshTestTransport func(*http.Request) (*http.Response, error)

func (f refreshTestTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestKMSRefreshResultTypes(t *testing.T) {
	for _, tc := range []struct {
		name    string
		body    string
		refresh func(*gophercloud.ServiceClient) retry.StateRefreshFunc
		check   func(*testing.T, interface{})
	}{
		{"key", `{"data":{"name":"test-key","type":"aes256-gcm96","deletion_allowed":true}}`,
			func(c *gophercloud.ServiceClient) retry.StateRefreshFunc {
				return keyStateRefreshFunc(c, "test-key", nil)
			},
			func(t *testing.T, result interface{}) {
				key, ok := result.(keys.GetKey)
				require.True(t, ok, "expected keys.GetKey, got %T", result)
				require.Equal(t, "test-key", key.Data.Name)
				require.True(t, key.Data.DeletionAllowed)
			}},
		{"secret", `{"data":{"data":{"value":"test"},"metadata":{"version":3}}}`,
			func(c *gophercloud.ServiceClient) retry.StateRefreshFunc {
				return secretStateRefreshFunc(c, "test-secret")
			},
			func(t *testing.T, result interface{}) {
				secret, ok := result.(secrets.GetSecret)
				require.True(t, ok, "expected secrets.GetSecret, got %T", result)
				require.Equal(t, "test", secret.Data.Data["value"])
				require.Equal(t, 3, secret.Data.Metadata.Version)
			}},
		{"protection true", `{"delete_protection":true}`,
			func(c *gophercloud.ServiceClient) retry.StateRefreshFunc {
				return secretDeleteProtectionStateRefreshFunc(c, "test-secret", nil)
			},
			func(t *testing.T, result interface{}) {
				protection, ok := result.(bool)
				require.True(t, ok, "expected bool, got %T", result)
				require.True(t, protection)
			}},
		{"protection false", `{"delete_protection":false}`,
			func(c *gophercloud.ServiceClient) retry.StateRefreshFunc {
				return secretDeleteProtectionStateRefreshFunc(c, "test-secret", nil)
			},
			func(t *testing.T, result interface{}) {
				protection, ok := result.(bool)
				require.True(t, ok, "expected bool, got %T", result)
				require.False(t, protection)
			}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				calls := 0
				provider := &gophercloud.ProviderClient{}
				provider.HTTPClient.Transport = refreshTestTransport(func(r *http.Request) (*http.Response, error) {
					calls++
					status, body := http.StatusOK, tc.body
					if calls == 1 {
						status, body = http.StatusBadRequest, `{"errors":["unfinished operation"]}`
					}
					return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
				})
				client := &gophercloud.ServiceClient{ProviderClient: provider, Endpoint: "https://kms.example.test/"}
				result, err := kmsStateChangeConf(tc.refresh(client), kmsStateReady).WaitForStateContext(context.Background())
				require.NoError(t, err)
				tc.check(t, result)
				require.Equal(t, 2, calls)
			})
		})
	}
}
