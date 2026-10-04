package kms_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kms"
)

func TestKMSSecretResourceRead(t *testing.T) {
	body := `{"data":{"data":{"username":"alice","password":"test-password"},"metadata":{"created_time":"2026-09-01T10:00:00Z","version":3}}}`
	config := newSecretDataSourceTestConfig(t, http.StatusOK, body)
	ds := kms.ResourceSecret()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"path": "test-secret", "data": map[string]interface{}{"stale": "value"},
		"data_json": `{"stale":"value"}`, "created_time": "stale-time", "version": 1,
	})
	d.SetId("test-secret")

	diags := ds.ReadContext(context.Background(), d, config)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "test-secret", d.Id())
	assert.Equal(t, "test-secret", d.Get("path"))
	assert.Equal(t, map[string]interface{}{"username": "alice", "password": "test-password"}, d.Get("data"))
	assert.JSONEq(t, `{"username":"alice","password":"test-password"}`, d.Get("data_json").(string))
	assert.Equal(t, "2026-09-01 10:00:00 +0000 UTC", d.Get("created_time"))
	assert.Equal(t, 3, d.Get("version"))
}

func TestKMSSecretResourceRead_errors(t *testing.T) {
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
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newSecretDataSourceTestConfig(t, tc.status, tc.body)
			ds := kms.ResourceSecret()
			d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"path": "test-secret"})
			d.SetId("test-secret")

			diags := ds.ReadContext(context.Background(), d, config)
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary, "Error getting secret:")
			assert.Equal(t, "test-secret", d.Id())
		})
	}
}

func TestKMSSecretResourceRead_clientError(t *testing.T) {
	ds := kms.ResourceSecret()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"path": "test-secret"})
	d.SetId("test-secret")
	config := &secretDataSourceTestConfig{err: errors.New("client unavailable")}

	diags := ds.ReadContext(context.Background(), d, config)
	require.True(t, diags.HasError())
	require.Len(t, diags, 1)
	assert.Equal(t, "Error creating VKCS KMS client: client unavailable", diags[0].Summary)
	assert.Equal(t, "test-secret", d.Id())
}
