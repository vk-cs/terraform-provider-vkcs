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

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"
)

func TestAccKMSSecretsListDataSource_basic(t *testing.T) {
	name := sdkacctest.RandomWithPrefix("tf-acc-kms-secret")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKMSSecretsListDataSourceBasic, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secrets_list.secrets", kms.EntityFieldID),
					resource.TestCheckResourceAttrSet("data.vkcs_kms_secrets_list.secrets", kms.SecretsListFieldSecretsCount),
					resource.TestCheckTypeSetElemAttr("data.vkcs_kms_secrets_list.secrets", kms.SecretsListFieldSecrets+".*", name),
				),
			},
		},
	})
}

const testAccKMSSecretsListDataSourceBasic = `
resource "vkcs_kms_secret" "secret" {
  path = %q
  delete_protection = false
  data_json = jsonencode({ password = "test-password" })
}

data "vkcs_kms_secrets_list" "secrets" {
  depends_on = [vkcs_kms_secret.secret]
}
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
			ds := newFrameworkTestEndpoint(t, kms.NewSecretsListDataSource())
			d := newFrameworkTestData(t, ds, map[string]interface{}{
				kms.SecretsListFieldSecrets: []interface{}{"stale-secret"}, kms.SecretsListFieldSecretsCount: 1,
			})

			diags := ds.ReadContext(context.Background(), d, config)
			require.False(t, diags.HasError(), "%v", diags)
			assert.Equal(t, tc.wantID, d.Id())
			assert.Equal(t, len(tc.wantSecrets), d.Get(kms.SecretsListFieldSecretsCount))
			assert.Equal(t, tc.wantSecrets, d.Get(kms.SecretsListFieldSecrets))
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
		{name: "invalid keys type", status: http.StatusOK, body: `{"data":{"keys":123}}`},
		{name: "invalid key element", status: http.StatusOK, body: `{"data":{"keys":[123]}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newSecretsListDataSourceTestConfig(t, tc.status, tc.body)
			ds := newFrameworkTestEndpoint(t, kms.NewSecretsListDataSource())
			d := newFrameworkTestData(t, ds, nil)

			diags := ds.ReadContext(context.Background(), d, config)
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary(), "Error listing secrets")
			assert.Empty(t, d.Id())
		})
	}
}

func TestKMSSecretsListDataSourceRead_clientError(t *testing.T) {
	ds := newFrameworkTestEndpoint(t, kms.NewSecretsListDataSource())
	d := newFrameworkTestData(t, ds, nil)
	config := &secretDataSourceTestConfig{err: errors.New("client unavailable")}

	diags := ds.ReadContext(context.Background(), d, config)
	require.True(t, diags.HasError())
	require.Len(t, diags, 1)
	assert.Equal(t, "Error creating VKCS KMS client", diags[0].Summary())
	assert.Equal(t, "client unavailable", diags[0].Detail())
	assert.Empty(t, d.Id())
}

func newSecretsListDataSourceTestConfig(t *testing.T, status int, body string) *secretDataSourceTestConfig {
	t.Helper()
	calls := 0
	t.Cleanup(func() { assert.Equal(t, 1, calls, "expected one KMS request") })
	provider := &gophercloud.ProviderClient{}
	provider.HTTPClient.Transport = secretDataSourceTestTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		assert.Empty(t, r.URL.Fragment)
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
