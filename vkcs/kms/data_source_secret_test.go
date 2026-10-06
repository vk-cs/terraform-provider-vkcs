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
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"
)

func TestAccKMSSecretDataSource_basic(t *testing.T) {
	name := sdkacctest.RandomWithPrefix("tf-acc-kms-secret")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKMSSecretDataSourceBasic, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_kms_secret.secret", "id", name),
					resource.TestCheckResourceAttr("data.vkcs_kms_secret.secret", "path", name),
					resource.TestCheckResourceAttrPair("data.vkcs_kms_secret.secret", "data_json", "vkcs_kms_secret.secret", "data_json"),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secret.secret", "created_time"),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secret.secret", "version"),
				),
			},
		},
	})
}

const testAccKMSSecretDataSourceBasic = `
resource "vkcs_kms_secret" "secret" {
  path = %q
  data_json = jsonencode({ password = "test-password" })
}

data "vkcs_kms_secret" "secret" {
  path = vkcs_kms_secret.secret.path
}
`

func TestKMSSecretDataSourceRead(t *testing.T) {
	body := `{"data":{"data":{"username":"alice","password":"test-password"},"metadata":{"created_time":"2026-09-01T10:00:00Z","version":3}}}`
	config := newSecretDataSourceTestConfig(t, http.StatusOK, body)
	ds := kms.DataSourceSecret()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"path": "test-secret", "data_json": `{"stale":"value"}`, "created_time": "stale-time", "version": 1})

	d.SetId("stale-id")
	diags := ds.ReadContext(context.Background(), d, config)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "test-secret", d.Id())
	assert.Equal(t, "test-secret", d.Get("path"))
	assert.JSONEq(t, `{"username":"alice","password":"test-password"}`, d.Get("data_json").(string))
	assert.Equal(t, "2026-09-01 10:00:00 +0000 UTC", d.Get("created_time"))
	assert.Equal(t, 3, d.Get("version"))
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
		{name: "invalid data type", status: http.StatusOK, body: `{"data":{"data":{"count":123}}}`},
		{name: "invalid timestamp", status: http.StatusOK, body: `{"data":{"metadata":{"created_time":"invalid"}}}`},
		{name: "invalid version", status: http.StatusOK, body: `{"data":{"metadata":{"version":"invalid"}}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newSecretDataSourceTestConfig(t, tc.status, tc.body)
			ds := kms.DataSourceSecret()
			d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"path": "test-secret"})

			diags := ds.ReadContext(context.Background(), d, config)
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary, "Error getting secret:")
			assert.Empty(t, d.Id())
		})
	}
}

func TestKMSSecretDataSourceRead_clientError(t *testing.T) {
	ds := kms.DataSourceSecret()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"path": "test-secret"})
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
	calls := 0
	t.Cleanup(func() { assert.Equal(t, 1, calls, "expected one KMS request") })
	provider := &gophercloud.ProviderClient{}
	provider.HTTPClient.Transport = secretDataSourceTestTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		assert.Empty(t, r.URL.Fragment)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/kms/user/v1/secret/data/test-secret", r.URL.EscapedPath())
		assert.Empty(t, r.URL.RawQuery)
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

func TestKMSSecretDataSourceRead_escapedPath(t *testing.T) {
	ds := kms.DataSourceSecret()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"path": "secret /?#%"})
	config := newKMSResourceTestConfig(t, kmsResourceRequest{
		method:   http.MethodGet,
		path:     "secret/data/secret%20%2F%3F%23%25",
		status:   http.StatusOK,
		response: `{"data":{"data":{"password":"test-password"},"metadata":{"created_time":"2026-09-01T10:00:00Z","version":1}}}`,
	})
	diags := ds.ReadContext(context.Background(), d, config)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "secret /?#%", d.Id())
	assert.Equal(t, "secret /?#%", d.Get("path"))
	assert.JSONEq(t, `{"password":"test-password"}`, d.Get("data_json").(string))
}
