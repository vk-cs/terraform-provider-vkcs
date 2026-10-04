package kms_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"
)

func TestKMSKeyResourceRead(t *testing.T) {
	for _, tc := range []struct {
		name            string
		path            string
		deletionAllowed bool
	}{
		{name: "test-key", path: "/kms/user/v1/transit/keys/test-key", deletionAllowed: true},
		{name: "key /?#%", path: "/kms/user/v1/transit/keys/key%20%2F%3F%23%25", deletionAllowed: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := fmt.Sprintf(`{"data":{"name":%q,"type":"aes256-gcm96","deletion_allowed":%t}}`, tc.name, tc.deletionAllowed)
			config := newKeyDataSourceTestConfig(t, tc.name, tc.path, http.StatusOK, body)
			ds := kms.ResourceKey()
			d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
				"name": tc.name, "type": "stale-type", "deletion_allowed": !tc.deletionAllowed,
			})
			d.SetId("stale-id")

			diags := ds.ReadContext(context.Background(), d, config)
			require.False(t, diags.HasError(), "%v", diags)
			assert.Equal(t, tc.name, d.Id())
			assert.Equal(t, tc.name, d.Get("name"))
			assert.Equal(t, "aes256-gcm96", d.Get("type"))
			assert.Equal(t, tc.deletionAllowed, d.Get("deletion_allowed"))
		})
	}
}

func TestKMSKeyResourceRead_errors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{name: "not found", status: http.StatusNotFound, body: `{"errors":["key not found"]}`},
		{name: "forbidden", status: http.StatusForbidden, body: `{"errors":["permission denied"]}`},
		{name: "invalid JSON", status: http.StatusOK, body: `invalid`},
		{name: "invalid field type", status: http.StatusOK, body: `{"data":{"deletion_allowed":"yes"}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newKeyDataSourceTestConfig(t, "test-key", "/kms/user/v1/transit/keys/test-key", tc.status, tc.body)
			ds := kms.ResourceKey()
			d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"name": "test-key"})
			d.SetId("test-key")

			diags := ds.ReadContext(context.Background(), d, config)
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary, "Error listing secrets:")
			assert.Equal(t, "test-key", d.Id())
		})
	}
}

func TestKMSKeyResourceRead_clientError(t *testing.T) {
	ds := kms.ResourceKey()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"name": "test-key"})
	d.SetId("test-key")
	config := &secretDataSourceTestConfig{err: errors.New("client unavailable")}

	diags := ds.ReadContext(context.Background(), d, config)
	require.True(t, diags.HasError())
	require.Len(t, diags, 1)
	assert.Equal(t, "Error creating VKCS KMS client: client unavailable", diags[0].Summary)
	assert.Equal(t, "test-key", d.Id())
}
