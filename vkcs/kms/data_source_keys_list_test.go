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

func TestKMSKeysListDataSourceRead(t *testing.T) {
	for _, tc := range []struct {
		name     string
		body     string
		wantKeys []interface{}
		wantID   string
	}{
		{
			name:     "multiple keys",
			body:     `{"data":{"keys":["test-key","another-key"]}}`,
			wantKeys: []interface{}{"test-key", "another-key"},
			wantID:   "eb2e78d667d7f3ed30a8c3a93174895e2fc91e317acdbb28d7f53250114bfc07",
		},
		{
			name:     "empty list",
			body:     `{"data":{"keys":[]}}`,
			wantKeys: []interface{}{},
			wantID:   "19e89348f2a9d5f3d0c5fca8e2a7068d9c5d71687a355f759009a5ed2527eb2c",
		},
		{
			name:     "null list",
			body:     `{"data":{"keys":null}}`,
			wantKeys: []interface{}{},
			wantID:   "19e89348f2a9d5f3d0c5fca8e2a7068d9c5d71687a355f759009a5ed2527eb2c",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newKeysListDataSourceTestConfig(t, http.StatusOK, tc.body)
			ds := kms.DataSourceKeysList()
			d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
				"keys": []interface{}{"stale-key"}, "keys_count": 1,
			})

			diags := ds.ReadContext(context.Background(), d, config)
			require.False(t, diags.HasError(), "%v", diags)
			assert.Equal(t, tc.wantID, d.Id())
			assert.Equal(t, len(tc.wantKeys), d.Get("keys_count"))
			assert.Equal(t, tc.wantKeys, d.Get("keys"))
		})
	}
}

func TestKMSKeysListDataSourceRead_errors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{name: "not found", status: http.StatusNotFound, body: `{"errors":["keys not found"]}`},
		{name: "forbidden", status: http.StatusForbidden, body: `{"errors":["permission denied"]}`},
		{name: "invalid JSON", status: http.StatusOK, body: `invalid`},
		{name: "invalid keys type", status: http.StatusOK, body: `{"data":{"keys":123}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newKeysListDataSourceTestConfig(t, tc.status, tc.body)
			ds := kms.DataSourceKeysList()
			d := schema.TestResourceDataRaw(t, ds.Schema, nil)

			diags := ds.ReadContext(context.Background(), d, config)
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary, "Error listing keys:")
			assert.Empty(t, d.Id())
		})
	}
}

func TestKMSKeysListDataSourceRead_clientError(t *testing.T) {
	ds := kms.DataSourceKeysList()
	d := schema.TestResourceDataRaw(t, ds.Schema, nil)
	config := &secretDataSourceTestConfig{err: errors.New("client unavailable")}

	diags := ds.ReadContext(context.Background(), d, config)
	require.True(t, diags.HasError())
	require.Len(t, diags, 1)
	assert.Equal(t, "Error creating VKCS KMS client: client unavailable", diags[0].Summary)
	assert.Empty(t, d.Id())
}

func newKeysListDataSourceTestConfig(t *testing.T, status int, body string) *secretDataSourceTestConfig {
	t.Helper()
	calls := 0
	t.Cleanup(func() { assert.Equal(t, 1, calls) })
	provider := &gophercloud.ProviderClient{}
	provider.HTTPClient.Transport = secretDataSourceTestTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/kms/user/v1/transit/keys", r.URL.Path)
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
