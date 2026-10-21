package kms_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"testing/synctest"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
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
			ds := kms.NewKeyResource()
			var schema fwresource.SchemaResponse
			ds.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{
				kms.KeyFieldName: tc.name, kms.KeyFieldType: "stale-type", kms.KeyFieldDeletionAllowed: !tc.deletionAllowed,
			})
			require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "stale-id").HasError())

			ds.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
			response := fwresource.ReadResponse{State: d}
			ds.Read(context.Background(), fwresource.ReadRequest{State: d}, &response)
			d = response.State
			diags := response.Diagnostics
			require.False(t, diags.HasError(), "%v", diags)
			var attribute1 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute1).HasError())
			assert.Equal(t, "stale-id", attribute1.ValueString())
			var attribute2 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeyFieldName), &attribute2).HasError())
			assert.Equal(t, tc.name, attribute2.ValueString())
			var attribute3 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeyFieldType), &attribute3).HasError())
			assert.Equal(t, "aes256-gcm96", attribute3.ValueString())
			var attribute4 types.Bool
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeyFieldDeletionAllowed), &attribute4).HasError())
			assert.Equal(t, tc.deletionAllowed, attribute4.ValueBool())
		})
	}
}

func TestKMSKeyResourceRead_errors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{name: "forbidden", status: http.StatusForbidden, body: `{"errors":["permission denied"]}`},
		{name: "invalid JSON", status: http.StatusOK, body: `invalid`},
		{name: "invalid field type", status: http.StatusOK, body: `{"data":{"deletion_allowed":"yes"}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newKeyDataSourceTestConfig(t, "test-key", "/kms/user/v1/transit/keys/test-key", tc.status, tc.body)
			ds := kms.NewKeyResource()
			var schema fwresource.SchemaResponse
			ds.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.KeyFieldType: "aes256-gcm96", kms.KeyFieldDeletionAllowed: false, kms.KeyFieldName: "test-key"})
			require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "test-key").HasError())

			ds.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
			response := fwresource.ReadResponse{State: d}
			ds.Read(context.Background(), fwresource.ReadRequest{State: d}, &response)
			d = response.State
			diags := response.Diagnostics
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary(), "Error getting key")
			var attribute5 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute5).HasError())
			assert.Equal(t, "test-key", attribute5.ValueString())
		})
	}
}

func TestKMSKeyResourceRead_clientError(t *testing.T) {
	ds := kms.NewKeyResource()
	var schema fwresource.SchemaResponse
	ds.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.KeyFieldType: "aes256-gcm96", kms.KeyFieldDeletionAllowed: false, kms.KeyFieldName: "test-key"})
	require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "test-key").HasError())
	config := &secretDataSourceTestConfig{err: errors.New("client unavailable")}

	ds.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
	response := fwresource.ReadResponse{State: d}
	ds.Read(context.Background(), fwresource.ReadRequest{State: d}, &response)
	d = response.State
	diags := response.Diagnostics
	require.True(t, diags.HasError())
	require.Len(t, diags, 1)
	assert.Equal(t, "Error creating VKCS KMS client", diags[0].Summary())
	assert.Equal(t, "client unavailable", diags[0].Detail())
	var attribute6 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute6).HasError())
	assert.Equal(t, "test-key", attribute6.ValueString())
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
			r := kms.NewKeyResource()
			var schema fwresource.SchemaResponse
			r.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
			input := map[string]interface{}{kms.KeyFieldType: "aes256-gcm96", kms.KeyFieldDeletionAllowed: false, kms.KeyFieldName: "test-key"}
			if tc.allowed {
				input[kms.KeyFieldType] = tc.keyType
				input[kms.KeyFieldDeletionAllowed] = true
			}
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, input)
			config := newKMSResourceTestConfig(
				t,
				kmsResourceRequest{"POST", "transit/keys/test-key", fmt.Sprintf(`{"type":%q}`, tc.keyType), 200, `{"data":{"name":"test-key"}}`},
				kmsResourceRequest{"POST", "transit/keys/test-key/config", fmt.Sprintf(`{"deletion_allowed":%t}`, tc.allowed), 200, `{"data":{}}`},
				kmsResourceRequest{"GET", "transit/keys/test-key", "", 200, fmt.Sprintf(`{"data":{"deletion_allowed":%t}}`, tc.allowed)},
			)
			r.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
			response := fwresource.CreateResponse{State: d}
			r.Create(context.Background(), fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: d.Schema, Raw: d.Raw}}, &response)
			d = response.State
			diags := response.Diagnostics
			require.False(t, diags.HasError(), "%v", diags)
			var attribute7 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute7).HasError())
			assert.Equal(t, "test-key", attribute7.ValueString())
			var attribute8 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeyFieldType), &attribute8).HasError())
			assert.Equal(t, tc.keyType, attribute8.ValueString())
			var attribute9 types.Bool
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.KeyFieldDeletionAllowed), &attribute9).HasError())
			assert.Equal(t, tc.allowed, attribute9.ValueBool())
		})
	}
}

func TestKMSKeyResourceUpdate(t *testing.T) {
	for _, allowed := range []bool{true, false} {
		t.Run(fmt.Sprint(allowed), func(t *testing.T) {
			r := kms.NewKeyResource()
			var schema fwresource.SchemaResponse
			r.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.KeyFieldType: "aes256-gcm96", kms.KeyFieldName: "test-key", kms.KeyFieldDeletionAllowed: allowed})
			require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "key /?#%").HasError())
			config := newKMSResourceTestConfig(t, kmsResourceRequest{"POST", "transit/keys/key%20%2F%3F%23%25/config", fmt.Sprintf(`{"deletion_allowed":%t}`, allowed), 200, `{"data":{}}`}, kmsResourceRequest{"GET", "transit/keys/key%20%2F%3F%23%25", "", 200, fmt.Sprintf(`{"data":{"deletion_allowed":%t}}`, allowed)})
			r.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
			response := fwresource.UpdateResponse{State: d}
			r.Update(context.Background(), fwresource.UpdateRequest{Plan: tfsdk.Plan{Schema: d.Schema, Raw: d.Raw}}, &response)
			d = response.State
			diags := response.Diagnostics
			require.False(t, diags.HasError(), "%v", diags)
			var attribute10 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute10).HasError())
			assert.Equal(t, "key /?#%", attribute10.ValueString())
		})
	}
}

func TestKMSKeyResourceDelete(t *testing.T) {
	r := kms.NewKeyResource()
	var schema fwresource.SchemaResponse
	r.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.KeyFieldType: "aes256-gcm96", kms.KeyFieldDeletionAllowed: false, kms.KeyFieldName: "key /?#%"})
	require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "key /?#%").HasError())
	config := newKMSResourceTestConfig(t, kmsResourceRequest{"DELETE", "transit/keys/key%20%2F%3F%23%25", "", 204, ""}, kmsResourceRequest{"GET", "transit/keys/key%20%2F%3F%23%25", "", 404, `{"errors":["key not found"]}`})
	r.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
	response := fwresource.DeleteResponse{State: d}
	r.Delete(context.Background(), fwresource.DeleteRequest{State: d}, &response)
	d = response.State
	diags := response.Diagnostics
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
				synctest.Test(t, func(t *testing.T) {
					r := kms.NewKeyResource()
					var schema fwresource.SchemaResponse
					r.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
					d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.KeyFieldType: "aes256-gcm96", kms.KeyFieldDeletionAllowed: false, kms.KeyFieldName: "test-key"})
					action := "Create"
					request := kmsResourceRequest{"POST", "transit/keys/test-key", `{"type":"aes256-gcm96"}`, failure.status, failure.body}
					prefix := "Error creating key"
					var requests []kmsResourceRequest
					switch operation {
					case "configure after create":
						requests = append(requests, kmsResourceRequest{"POST", "transit/keys/test-key", request.body, 200, `{"data":{}}`})
						fallthrough
					case "update":
						request.path = "transit/keys/test-key/config"
						request.body = `{"deletion_allowed":false}`
						prefix = "Error updating key configuration"
						if operation == "update" {
							action = "Update"
							require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "test-key").HasError())
						}
					case "delete":
						action = "Delete"
						require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "test-key").HasError())
						request.method = "DELETE"
						request.path = "transit/keys/test-key"
						request.body = ""
						prefix = "Error deleting key"
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
					r.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
					var diags diag.Diagnostics
					switch action {
					case "Create":
						response := fwresource.CreateResponse{State: d}
						r.Create(context.Background(), fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: d.Schema, Raw: d.Raw}}, &response)
						d = response.State
						diags = response.Diagnostics
					case "Update":
						response := fwresource.UpdateResponse{State: d}
						r.Update(context.Background(), fwresource.UpdateRequest{Plan: tfsdk.Plan{Schema: d.Schema, Raw: d.Raw}}, &response)
						d = response.State
						diags = response.Diagnostics
					case "Delete":
						response := fwresource.DeleteResponse{State: d}
						r.Delete(context.Background(), fwresource.DeleteRequest{State: d}, &response)
						d = response.State
						diags = response.Diagnostics
					}
					require.True(t, diags.HasError())
					require.Len(t, diags, 1)
					assert.Contains(t, diags[0].Summary(), prefix)
					if operation == "create" || operation == "configure after create" {
						var attribute11 types.String
						require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute11).HasError())
						assert.Empty(t, attribute11.ValueString())
					} else {
						var attribute12 types.String
						require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute12).HasError())
						assert.Equal(t, "test-key", attribute12.ValueString())
					}
				})
			})
		}
	}
}

func TestKMSKeyResourceCreate_waitsForConfiguration(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		r := kms.NewKeyResource()
		var schema fwresource.SchemaResponse
		r.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
		d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.KeyFieldType: "aes256-gcm96", kms.KeyFieldName: "test-key", kms.KeyFieldDeletionAllowed: true})
		config := newKMSResourceTestConfig(t,
			kmsResourceRequest{"POST", "transit/keys/test-key", `{"type":"aes256-gcm96"}`, 200, `{"data":{}}`},
			kmsResourceRequest{"POST", "transit/keys/test-key/config", `{"deletion_allowed":true}`, 200, `{"data":{}}`},
			kmsResourceRequest{"GET", "transit/keys/test-key", "", 200, `{"data":{"deletion_allowed":false}}`},
			kmsResourceRequest{"GET", "transit/keys/test-key", "", 200, `{"data":{"deletion_allowed":true}}`},
		)
		r.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
		response := fwresource.CreateResponse{State: d}
		r.Create(context.Background(), fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: d.Schema, Raw: d.Raw}}, &response)
		d = response.State
		diags := response.Diagnostics
		require.False(t, diags.HasError(), "%v", diags)
	})
}

func TestKMSKeyResourceDelete_waitsForOperationAndDeletion(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		r := kms.NewKeyResource()
		var schema fwresource.SchemaResponse
		r.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
		d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.KeyFieldType: "aes256-gcm96", kms.KeyFieldDeletionAllowed: false, kms.KeyFieldName: "test-key"})
		require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "test-key").HasError())
		config := newKMSResourceTestConfig(t,
			kmsResourceRequest{"DELETE", "transit/keys/test-key", "", 400, `{"errors":["key test-key has an unfinished operation"]}`},
			kmsResourceRequest{"DELETE", "transit/keys/test-key", "", 204, ""},
			kmsResourceRequest{"GET", "transit/keys/test-key", "", 200, `{"data":{"name":"test-key"}}`},
			kmsResourceRequest{"GET", "transit/keys/test-key", "", 404, `{"errors":["key not found"]}`},
		)
		r.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
		response := fwresource.DeleteResponse{State: d}
		r.Delete(context.Background(), fwresource.DeleteRequest{State: d}, &response)
		d = response.State
		diags := response.Diagnostics
		require.False(t, diags.HasError(), "%v", diags)
	})
}

func TestKMSKeyResourceDelete_pollError(t *testing.T) {
	r := kms.NewKeyResource()
	var schema fwresource.SchemaResponse
	r.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.KeyFieldType: "aes256-gcm96", kms.KeyFieldDeletionAllowed: false, kms.KeyFieldName: "test-key"})
	require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "test-key").HasError())
	config := newKMSResourceTestConfig(t,
		kmsResourceRequest{"DELETE", "transit/keys/test-key", "", 204, ""},
		kmsResourceRequest{"GET", "transit/keys/test-key", "", 403, `{"errors":["permission denied"]}`},
	)
	r.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
	response := fwresource.DeleteResponse{State: d}
	r.Delete(context.Background(), fwresource.DeleteRequest{State: d}, &response)
	d = response.State
	diags := response.Diagnostics
	require.True(t, diags.HasError())
	assert.Equal(t, "Error waiting for key deletion", diags[0].Summary())
	assert.Contains(t, diags[0].Detail(), "permission denied")
}

func TestAccKMSKeyResource_basic(t *testing.T) {
	name := "tfacc-kms-key-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccKMSPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{Config: fmt.Sprintf(testAccKMSKeyResourceBasic, name), Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("vkcs_kms_key.key", "id", name),
				resource.TestCheckResourceAttr("vkcs_kms_key.key", "type", "aes256-gcm96"),
				resource.TestCheckResourceAttr("vkcs_kms_key.key", "deletion_allowed", "false"),
			)},
			{Config: fmt.Sprintf(testAccKMSKeyResourceConfig, name, "aes256-gcm96", true), Check: resource.TestCheckResourceAttr("vkcs_kms_key.key", "deletion_allowed", "true")},
			{Config: fmt.Sprintf(testAccKMSKeyResourceConfig, name+"-replacement", "aes128-gcm96", true), Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("vkcs_kms_key.key", "id", name+"-replacement"),
				resource.TestCheckResourceAttr("vkcs_kms_key.key", "type", "aes128-gcm96"),
			)},
		},
	})
}

func TestAccKMSKeyResource_keyTypes(t *testing.T) {
	for _, keyType := range []string{"aes128-gcm96", "aes256-gcm96", "chacha20-poly1305", "xchacha20-poly1305"} {
		t.Run(keyType, func(t *testing.T) {
			name := "tfacc-kms-key-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccKMSPreCheck(t) },
				ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
				CheckDestroy:             testAccKMSCheckDestroy,
				Steps:                    []resource.TestStep{{Config: fmt.Sprintf(testAccKMSKeyResourceConfig, name, keyType, true), Check: resource.TestCheckResourceAttr("vkcs_kms_key.key", "type", keyType)}},
			})
		})
	}
}

const testAccKMSKeyResourceBasic = `
resource "vkcs_kms_key" "key" {
  name = %q
}
`

const testAccKMSKeyResourceConfig = `
resource "vkcs_kms_key" "key" {
  name             = %q
  type             = %q
  deletion_allowed = %t
}
`
