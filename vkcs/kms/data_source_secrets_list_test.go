package kms_test

import (
	"context"
	"errors"
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
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"
)

// The KMS resource does not implement creation yet, so this read-only acceptance
// test uses an existing secret supplied through OS_KMS_SECRET_NAME.
func TestAccKMSSecretsListDataSource_basic(t *testing.T) {
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
				Config: testAccKMSSecretsListDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secrets_list.secrets", "id"),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secrets_list.secrets", "secrets_count"),
					resource.TestCheckTypeSetElemAttr("data.vkcs_kms_secrets_list.secrets", "secrets.*", name),
				),
			},
		},
	})
}

const testAccKMSSecretsListDataSourceBasic = `
data "vkcs_kms_secrets_list" "secrets" {}
`

func TestKMSSecretsListDataSourceRead(t *testing.T) {
	for _, tc := range []struct {
		name        string
		body        string
		wantSecrets []interface{}
		wantID      string
	}{
		{
			name:        "multiple secrets",
			body:        `{"data":{"keys":["test-secret","another-secret"]}}`,
			wantSecrets: []interface{}{"test-secret", "another-secret"},
			wantID:      "7a9fbe1cce12c710a5bdc4a5e3bd05de40d59ddf96b6104ea0376bc5f6fe485f",
		},
		{
			name:        "empty list",
			body:        `{"data":{"keys":[]}}`,
			wantSecrets: []interface{}{},
			wantID:      "19e89348f2a9d5f3d0c5fca8e2a7068d9c5d71687a355f759009a5ed2527eb2c",
		},
		{
			name:        "null list",
			body:        `{"data":{"keys":null}}`,
			wantSecrets: []interface{}{},
			wantID:      "19e89348f2a9d5f3d0c5fca8e2a7068d9c5d71687a355f759009a5ed2527eb2c",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newSecretsListDataSourceTestConfig(t, http.StatusOK, tc.body)
			ds := kms.DataSourceSecretsList()
			d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
				"secrets": []interface{}{"stale-secret"}, "secrets_count": 1,
			})

			diags := ds.ReadContext(context.Background(), d, config)
			require.False(t, diags.HasError(), "%v", diags)
			assert.Equal(t, tc.wantID, d.Id())
			assert.Equal(t, len(tc.wantSecrets), d.Get("secrets_count"))
			assert.Equal(t, tc.wantSecrets, d.Get("secrets"))
		})
	}
}

func TestKMSSecretsListDataSourceRead_errors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{name: "not found", status: http.StatusNotFound, body: `{"errors":["secrets not found"]}`},
		{name: "forbidden", status: http.StatusForbidden, body: `{"errors":["permission denied"]}`},
		{name: "invalid JSON", status: http.StatusOK, body: `invalid`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newSecretsListDataSourceTestConfig(t, tc.status, tc.body)
			ds := kms.DataSourceSecretsList()
			d := schema.TestResourceDataRaw(t, ds.Schema, nil)

			diags := ds.ReadContext(context.Background(), d, config)
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary, "Error listing secrets:")
			assert.Empty(t, d.Id())
		})
	}
}

func TestKMSSecretsListDataSourceRead_clientError(t *testing.T) {
	ds := kms.DataSourceSecretsList()
	d := schema.TestResourceDataRaw(t, ds.Schema, nil)
	config := &secretDataSourceTestConfig{err: errors.New("client unavailable")}

	diags := ds.ReadContext(context.Background(), d, config)
	require.True(t, diags.HasError())
	require.Len(t, diags, 1)
	assert.Equal(t, "Error creating VKCS KMS client: client unavailable", diags[0].Summary)
	assert.Empty(t, d.Id())
}

func newSecretsListDataSourceTestConfig(t *testing.T, status int, body string) *secretDataSourceTestConfig {
	t.Helper()
	provider := &gophercloud.ProviderClient{}
	provider.HTTPClient.Transport = secretDataSourceTestTransport(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/kms/user/v1/secret/metadata/", r.URL.Path)
		assert.Equal(t, "list=true", r.URL.RawQuery)
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
