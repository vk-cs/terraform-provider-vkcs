package kms_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"
)

// The KMS resource does not implement creation yet, so this read-only acceptance
// test uses an existing secret supplied through OS_KMS_SECRET_NAME.
func TestAccKMSSecretDataSource_basic(t *testing.T) {
	name := os.Getenv("OS_KMS_SECRET_NAME")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if name == "" {
				t.Fatal("OS_KMS_SECRET_NAME must be set for this acceptance test")
			}
		},
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKMSSecretDataSourceBasic, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_kms_secret.secret", "id", name),
					resource.TestCheckResourceAttr("data.vkcs_kms_secret.secret", "name", name),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secret.secret", "data_json"),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secret.secret", "created_time"),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secret.secret", "version"),
				),
			},
		},
	})
}

const testAccKMSSecretDataSourceBasic = `
data "vkcs_kms_secret" "secret" {
  name = %q
}
`

func TestKMSSecretDataSourceRead(t *testing.T) {
	for _, tc := range []struct {
		name         string
		deletionJSON string
		wantDeletion string
	}{
		{name: "active", deletionJSON: "null"},
		{
			name:         "scheduled deletion",
			deletionJSON: `"2026-10-01T12:00:00Z"`,
			wantDeletion: "2026-10-01 12:00:00 +0000 UTC",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := fmt.Sprintf(`{"data":{"data":{"username":"alice","password":"test-password"},"metadata":{"created_time":"2026-09-01T10:00:00Z","deletion_time":%s,"version":3}}}`, tc.deletionJSON)
			config := newSecretDataSourceTestConfig(t, http.StatusOK, body)
			ds := kms.DataSourceSecret()
			d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"name": "test-secret"})

			diags := ds.ReadContext(context.Background(), d, config)
			require.False(t, diags.HasError(), "%v", diags)
			assert.Equal(t, "test-secret", d.Id())
			assert.Equal(t, "test-secret", d.Get("name"))
			assert.Equal(t, map[string]interface{}{"username": "alice", "password": "test-password"}, d.Get("data"))
			assert.JSONEq(t, `{"username":"alice","password":"test-password"}`, d.Get("data_json").(string))
			assert.Equal(t, "2026-09-01 10:00:00 +0000 UTC", d.Get("created_time"))
			assert.Equal(t, tc.wantDeletion, d.Get("deletion_time"))
			assert.Equal(t, 3, d.Get("version"))
		})
	}
}

func TestKMSSecretDataSourceRead_errors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{name: "not found", status: http.StatusNotFound, body: `{"errors":["secret not found"]}`},
		{name: "forbidden", status: http.StatusForbidden, body: `{"errors":["permission denied"]}`},
		{name: "invalid JSON", status: http.StatusOK, body: `invalid`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newSecretDataSourceTestConfig(t, tc.status, tc.body)
			ds := kms.DataSourceSecret()
			d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"name": "test-secret"})

			diags := ds.ReadContext(context.Background(), d, config)
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary, "Error listing secrets:")
			assert.Empty(t, d.Id())
		})
	}
}

func TestKMSSecretDataSourceRead_clientError(t *testing.T) {
	ds := kms.DataSourceSecret()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"name": "test-secret"})
	config := &secretDataSourceTestConfig{err: errors.New("client unavailable")}

	diags := ds.ReadContext(context.Background(), d, config)
	require.True(t, diags.HasError())
	require.Len(t, diags, 1)
	assert.Equal(t, "Error creating VKCS KMS client: client unavailable", diags[0].Summary)
	assert.Empty(t, d.Id())
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

func newSecretDataSourceTestConfig(t *testing.T, status int, body string) *secretDataSourceTestConfig {
	t.Helper()
	provider := &gophercloud.ProviderClient{}
	provider.HTTPClient.Transport = secretDataSourceTestTransport(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/kms/user/v1/secret/data/test-secret", r.URL.Path)
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
