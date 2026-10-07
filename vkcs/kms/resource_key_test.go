package kms_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
			assert.Equal(t, "stale-id", d.Id())
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
			assert.Contains(t, diags[0].Summary, "Error getting key:")
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

func TestKMSKeyResourceCreate(t *testing.T) {
	for _, tc := range []struct {
		name, keyType string
		allowed       bool
	}{
		{name: "defaults", keyType: "aes256-gcm96"},
		{name: "explicit settings", keyType: "aes128-gcm96", allowed: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := kms.ResourceKey()
			input := map[string]interface{}{"name": "test-key"}
			if tc.allowed {
				input["type"] = tc.keyType
				input["deletion_allowed"] = true
			}
			d := schema.TestResourceDataRaw(t, r.Schema, input)
			config := newKMSResourceTestConfig(
				t,
				kmsResourceRequest{"POST", "transit/keys/test-key", fmt.Sprintf(`{"type":%q}`, tc.keyType), 200, `{"data":{"name":"test-key"}}`},
				kmsResourceRequest{"POST", "transit/keys/test-key/config", fmt.Sprintf(`{"deletion_allowed":%t}`, tc.allowed), 200, `{"data":{}}`},
			)
			diags := r.CreateContext(context.Background(), d, config)
			require.False(t, diags.HasError(), "%v", diags)
			assert.Equal(t, "test-key", d.Id())
			assert.Equal(t, tc.keyType, d.Get("type"))
			assert.Equal(t, tc.allowed, d.Get("deletion_allowed"))
		})
	}
}

func TestKMSKeyResourceUpdate(t *testing.T) {
	for _, allowed := range []bool{true, false} {
		t.Run(fmt.Sprint(allowed), func(t *testing.T) {
			r := kms.ResourceKey()
			d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"name": "test-key", "deletion_allowed": allowed})
			d.SetId("key /?#%")
			config := newKMSResourceTestConfig(t, kmsResourceRequest{"POST", "transit/keys/key%20%2F%3F%23%25/config", fmt.Sprintf(`{"deletion_allowed":%t}`, allowed), 200, `{"data":{}}`})
			diags := r.UpdateContext(context.Background(), d, config)
			require.False(t, diags.HasError(), "%v", diags)
			assert.Equal(t, "key /?#%", d.Id())
		})
	}
}

func TestKMSKeyResourceDelete(t *testing.T) {
	r := kms.ResourceKey()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"name": "key /?#%"})
	d.SetId("key /?#%")
	config := newKMSResourceTestConfig(t, kmsResourceRequest{"DELETE", "transit/keys/key%20%2F%3F%23%25", "", 204, ""}, kmsResourceRequest{"GET", "transit/keys/key%20%2F%3F%23%25", "", 404, `{"errors":["key not found"]}`})
	diags := r.DeleteContext(context.Background(), d, config)
	require.False(t, diags.HasError(), "%v", diags)
}

func TestKMSKeyResourceWriteErrors(t *testing.T) {
	for _, operation := range []string{"create", "configure after create", "update", "delete"} {
		for _, failure := range []struct {
			name   string
			status int
			body   string
		}{
			{"forbidden", 403, `{"errors":["permission denied"]}`},
			{"invalid JSON", 200, `invalid`},
		} {
			if operation == "delete" && failure.name == "invalid JSON" {
				continue
			}
			t.Run(operation+"/"+failure.name, func(t *testing.T) {
				r := kms.ResourceKey()
				d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"name": "test-key"})
				var call func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics = r.CreateContext
				request := kmsResourceRequest{"POST", "transit/keys/test-key", `{"type":"aes256-gcm96"}`, failure.status, failure.body}
				prefix := "Error creating key:"
				var requests []kmsResourceRequest
				switch operation {
				case "configure after create":
					requests = append(requests, kmsResourceRequest{"POST", "transit/keys/test-key", request.body, 200, `{"data":{}}`})
					fallthrough
				case "update":
					request.path = "transit/keys/test-key/config"
					request.body = `{"deletion_allowed":false}`
					prefix = "Error updating key configuration:"
					if operation == "update" {
						call = r.UpdateContext
						d.SetId("test-key")
					}
				case "delete":
					call = r.DeleteContext
					d.SetId("test-key")
					request.method = "DELETE"
					request.path = "transit/keys/test-key"
					request.body = ""
					prefix = "Error deleting key:"
				}
				if operation == "create" && failure.name == "invalid JSON" {
					request.status = 200
				}
				requests = append(requests, request)
				if operation == "configure after create" {
					for range 9 {
						requests = append(requests, request)
					}
				}
				config := newKMSResourceTestConfig(t, requests...)
				diags := call(context.Background(), d, config)
				require.True(t, diags.HasError())
				require.Len(t, diags, 1)
				assert.Contains(t, diags[0].Summary, prefix)
				if operation == "create" || operation == "configure after create" {
					assert.Empty(t, d.Id())
				} else {
					assert.Equal(t, "test-key", d.Id())
				}
			})
		}
	}
}
