package kms_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		"path":      "test-secret",
		"data_json": `{"stale":"value"}`, "created_time": "stale-time", "version": 1,
	})
	d.SetId("test-secret")

	diags := ds.ReadContext(context.Background(), d, config)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "test-secret", d.Id())
	assert.Equal(t, "test-secret", d.Get("path"))
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

func TestKMSSecretResourceWrite(t *testing.T) {
	for _, operation := range []string{"create", "update"} {
		t.Run(operation, func(t *testing.T) {
			r := kms.ResourceSecret()
			d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"path": "secret /?#%", "data_json": `{"password":"new-value"}`, "created_time": "stale-time", "version": 1})
			var call func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics = r.CreateContext
			if operation == "update" {
				call = r.UpdateContext
				d.SetId("secret /?#%")
			}
			config := newKMSResourceTestConfig(t, kmsResourceRequest{"POST", "secret/data/secret%20%2F%3F%23%25", `{"data":"{\"password\":\"new-value\"}"}`, 200, `{"data":{"created_time":"2026-09-01T10:00:00Z","version":2}}`})
			diags := call(context.Background(), d, config)
			require.False(t, diags.HasError(), "%v", diags)
			assert.Equal(t, "secret /?#%", d.Id())
			assert.JSONEq(t, `{"password":"new-value"}`, d.Get("data_json").(string))
			assert.Equal(t, "2026-09-01 10:00:00 +0000 UTC", d.Get("created_time"))
			assert.Equal(t, 2, d.Get("version"))
		})
	}
}

func TestKMSSecretResourceDelete(t *testing.T) {
	r := kms.ResourceSecret()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"path": "secret /?#%"})
	d.SetId("secret /?#%")
	config := newKMSResourceTestConfig(t, kmsResourceRequest{"DELETE", "secret/metadata/secret%20%2F%3F%23%25", "", 204, ""})
	diags := r.DeleteContext(context.Background(), d, config)
	require.False(t, diags.HasError(), "%v", diags)
}

func TestKMSSecretResourceWriteErrors(t *testing.T) {
	for _, operation := range []string{"create", "update", "delete"} {
		for _, failure := range []struct {
			name   string
			status int
			body   string
		}{
			{"forbidden", 403, `{"errors":["permission denied"]}`},
			{"invalid JSON", 200, `invalid`},
			{"invalid timestamp", 200, `{"data":{"created_time":"invalid"}}`},
		} {
			if operation == "delete" && failure.status == 200 {
				continue
			}
			t.Run(operation+"/"+failure.name, func(t *testing.T) {
				r := kms.ResourceSecret()
				d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"path": "test-secret", "data_json": `{}`, "created_time": "old-time", "version": 1})
				var call func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics = r.CreateContext
				request := kmsResourceRequest{"POST", "secret/data/test-secret", `{"data":"{}"}`, failure.status, failure.body}
				prefix := "Error creating secret:"
				if operation != "create" {
					d.SetId("test-secret")
				}
				if operation == "update" {
					call = r.UpdateContext
				}
				if operation == "delete" {
					call = r.DeleteContext
					request.method = "DELETE"
					request.path = "secret/metadata/test-secret"
					request.body = ""
					prefix = "Error deleting secret:"
				}
				diags := call(context.Background(), d, newKMSResourceTestConfig(t, request))
				require.True(t, diags.HasError())
				require.Len(t, diags, 1)
				assert.Contains(t, diags[0].Summary, prefix)
				if operation == "create" {
					assert.Empty(t, d.Id())
				} else {
					assert.Equal(t, "test-secret", d.Id())
				}
				assert.Equal(t, "old-time", d.Get("created_time"))
				assert.Equal(t, 1, d.Get("version"))
			})
		}
	}
}
