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

func TestKMSSecretResourceRead(t *testing.T) {
	body := `{"data":{"data":{"username":"alice","password":"test-password"},"metadata":{"created_time":"2026-09-01T10:00:00Z","version":3}}}`
	config := newSecretDataSourceTestConfig(t, http.StatusOK, body, false)
	ds := kms.NewSecretResource()
	var schema fwresource.SchemaResponse
	ds.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.SecretFieldDeleteProtection: true,
		kms.SecretFieldPath:     "test-secret",
		kms.SecretFieldDataJSON: `{"stale":"value"}`, kms.SecretFieldCreatedTime: "stale-time", kms.SecretFieldVersion: 1,
	})
	require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "test-secret").HasError())

	ds.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
	response := fwresource.ReadResponse{State: d}
	ds.Read(context.Background(), fwresource.ReadRequest{State: d}, &response)
	d = response.State
	diags := response.Diagnostics
	require.False(t, diags.HasError(), "%v", diags)
	var attribute1 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute1).HasError())
	assert.Equal(t, "test-secret", attribute1.ValueString())
	var attribute2 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldPath), &attribute2).HasError())
	assert.Equal(t, "test-secret", attribute2.ValueString())
	var attribute3 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldDataJSON), &attribute3).HasError())
	assert.JSONEq(t, `{"username":"alice","password":"test-password"}`, attribute3.ValueString())
	var attribute4 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldCreatedTime), &attribute4).HasError())
	assert.Equal(t, "2026-09-01 10:00:00 +0000 UTC", attribute4.ValueString())
	var attribute5 types.Int64
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldVersion), &attribute5).HasError())
	assert.Equal(t, 3, int(attribute5.ValueInt64()))
	var attribute6 types.Bool
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldDeleteProtection), &attribute6).HasError())
	assert.Equal(t, false, attribute6.ValueBool())
}

func TestKMSSecretResourceRead_errors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{name: "forbidden", status: http.StatusForbidden, body: `{"errors":["permission denied"]}`},
		{name: "invalid JSON", status: http.StatusOK, body: `invalid`},
		{name: "invalid data type", status: http.StatusOK, body: `{"data":{"data":{"count":123}}}`},
		{name: "invalid timestamp", status: http.StatusOK, body: `{"data":{"metadata":{"created_time":"invalid"}}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := newSecretDataSourceTestConfig(t, tc.status, tc.body)
			ds := kms.NewSecretResource()
			var schema fwresource.SchemaResponse
			ds.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.SecretFieldDeleteProtection: true, kms.SecretFieldPath: "test-secret"})
			require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "test-secret").HasError())

			ds.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
			response := fwresource.ReadResponse{State: d}
			ds.Read(context.Background(), fwresource.ReadRequest{State: d}, &response)
			d = response.State
			diags := response.Diagnostics
			require.True(t, diags.HasError())
			require.Len(t, diags, 1)
			assert.Contains(t, diags[0].Summary(), "Error getting secret")
			var attribute7 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute7).HasError())
			assert.Equal(t, "test-secret", attribute7.ValueString())
		})
	}
}

func TestKMSSecretResourceRead_clientError(t *testing.T) {
	ds := kms.NewSecretResource()
	var schema fwresource.SchemaResponse
	ds.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.SecretFieldDeleteProtection: true, kms.SecretFieldPath: "test-secret"})
	require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "test-secret").HasError())
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
	var attribute8 types.String
	require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute8).HasError())
	assert.Equal(t, "test-secret", attribute8.ValueString())
}

func TestKMSSecretResourceWrite(t *testing.T) {
	for _, operation := range []string{"create", "update"} {
		t.Run(operation, func(t *testing.T) {
			r := kms.NewSecretResource()
			var schema fwresource.SchemaResponse
			r.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
			d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.SecretFieldDeleteProtection: true, kms.SecretFieldPath: "secret /?#%", kms.SecretFieldDataJSON: `{"password":"new-value"}`, kms.SecretFieldCreatedTime: "stale-time", kms.SecretFieldVersion: 1})
			action := "Create"
			if operation == "update" {
				action = "Update"
				require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "secret /?#%").HasError())
			}
			config := newKMSResourceTestConfig(t, kmsResourceRequest{"POST", "secret/data/secret%20%2F%3F%23%25", `{"data":{"password":"new-value"}}`, 200, `{"data":{"created_time":"2026-09-01T10:00:00Z","version":2}}`}, kmsResourceRequest{"PUT", "secret/metadata/delete-protection/secret%20%2F%3F%23%25", `{"delete_protection":true}`, 200, ""}, kmsResourceRequest{"GET", "secret/metadata/delete-protection/secret%20%2F%3F%23%25", "", 200, `{"delete_protection":true}`})
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
			require.False(t, diags.HasError(), "%v", diags)
			var attribute9 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute9).HasError())
			assert.Equal(t, "secret /?#%", attribute9.ValueString())
			var attribute10 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldDataJSON), &attribute10).HasError())
			assert.JSONEq(t, `{"password":"new-value"}`, attribute10.ValueString())
			var attribute11 types.String
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldCreatedTime), &attribute11).HasError())
			assert.Equal(t, "2026-09-01 10:00:00 +0000 UTC", attribute11.ValueString())
			var attribute12 types.Int64
			require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldVersion), &attribute12).HasError())
			assert.Equal(t, 2, int(attribute12.ValueInt64()))
		})
	}
}

func TestKMSSecretResourceDelete(t *testing.T) {
	r := kms.NewSecretResource()
	var schema fwresource.SchemaResponse
	r.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.SecretFieldDeleteProtection: true, kms.SecretFieldPath: "secret /?#%"})
	require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "secret /?#%").HasError())
	config := newKMSResourceTestConfig(t, kmsResourceRequest{"DELETE", "secret/metadata/secret%20%2F%3F%23%25", "", 204, ""}, kmsResourceRequest{"GET", "secret/data/secret%20%2F%3F%23%25", "", 404, `{"errors":["secret not found"]}`})
	r.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
	response := fwresource.DeleteResponse{State: d}
	r.Delete(context.Background(), fwresource.DeleteRequest{State: d}, &response)
	d = response.State
	diags := response.Diagnostics
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
				synctest.Test(t, func(t *testing.T) {
					r := kms.NewSecretResource()
					var schema fwresource.SchemaResponse
					r.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
					d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.SecretFieldDeleteProtection: true, kms.SecretFieldPath: "test-secret", kms.SecretFieldDataJSON: `{}`, kms.SecretFieldCreatedTime: "old-time", kms.SecretFieldVersion: 1})
					action := "Create"
					request := kmsResourceRequest{"POST", "secret/data/test-secret", `{"data":{}}`, failure.status, failure.body}
					prefix := "Error creating secret"
					if operation != "create" {
						require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "test-secret").HasError())
					}
					if operation == "update" {
						action = "Update"
					}
					if operation == "delete" {
						action = "Delete"
						request.method = "DELETE"
						request.path = "secret/metadata/test-secret"
						request.body = ""
						prefix = "Error deleting secret"
					}
					r.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: newKMSResourceTestConfig(t, request)}, &fwresource.ConfigureResponse{})
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
					if operation == "create" {
						var attribute13 types.String
						require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute13).HasError())
						assert.Empty(t, attribute13.ValueString())
					} else {
						var attribute14 types.String
						require.False(t, d.GetAttribute(context.Background(), path.Root(kms.EntityFieldID), &attribute14).HasError())
						assert.Equal(t, "test-secret", attribute14.ValueString())
					}
					var attribute15 types.String
					require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldCreatedTime), &attribute15).HasError())
					assert.Equal(t, "old-time", attribute15.ValueString())
					var attribute16 types.Int64
					require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldVersion), &attribute16).HasError())
					assert.Equal(t, 1, int(attribute16.ValueInt64()))
				})
			})
		}
	}
}

func TestKMSSecretResourceWrite_waitsForDeleteProtection(t *testing.T) {
	for _, operation := range []string{"create", "update"} {
		t.Run(operation, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				r := kms.NewSecretResource()
				var schema fwresource.SchemaResponse
				r.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
				d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.SecretFieldPath: "test-secret", kms.SecretFieldDataJSON: `{}`, kms.SecretFieldDeleteProtection: false})
				action := "Create"
				if operation == "update" {
					action = "Update"
					require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "test-secret").HasError())
				}
				config := newKMSResourceTestConfig(t,
					kmsResourceRequest{"POST", "secret/data/test-secret", `{"data":{}}`, 200, `{"data":{"created_time":"2026-09-01T10:00:00Z","version":2}}`},
					kmsResourceRequest{"PUT", "secret/metadata/delete-protection/test-secret", `{"delete_protection":false}`, 400, `{"errors":["secret has an unfinished operation"]}`},
					kmsResourceRequest{"PUT", "secret/metadata/delete-protection/test-secret", `{"delete_protection":false}`, 200, ""},
					kmsResourceRequest{"GET", "secret/metadata/delete-protection/test-secret", "", 200, `{"delete_protection":true}`},
					kmsResourceRequest{"GET", "secret/metadata/delete-protection/test-secret", "", 200, `{"delete_protection":false}`},
				)
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
				require.False(t, diags.HasError(), "%v", diags)
				var attribute17 types.Bool
				require.False(t, d.GetAttribute(context.Background(), path.Root(kms.SecretFieldDeleteProtection), &attribute17).HasError())
				assert.Equal(t, false, attribute17.ValueBool())
			})
		})
	}
}

func TestKMSSecretResourceDelete_waitsForOperationAndDeletion(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		r := kms.NewSecretResource()
		var schema fwresource.SchemaResponse
		r.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
		d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.SecretFieldDeleteProtection: true, kms.SecretFieldPath: "test-secret"})
		require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "test-secret").HasError())
		config := newKMSResourceTestConfig(t,
			kmsResourceRequest{"DELETE", "secret/metadata/test-secret", "", 400, `{"errors":["secret has an unfinished operation"]}`},
			kmsResourceRequest{"DELETE", "secret/metadata/test-secret", "", 204, ""},
			kmsResourceRequest{"GET", "secret/data/test-secret", "", 200, `{"data":{"data":{},"metadata":{"version":1}}}`},
			kmsResourceRequest{"GET", "secret/data/test-secret", "", 404, `{"errors":["secret not found"]}`},
		)
		r.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
		response := fwresource.DeleteResponse{State: d}
		r.Delete(context.Background(), fwresource.DeleteRequest{State: d}, &response)
		d = response.State
		diags := response.Diagnostics
		require.False(t, diags.HasError(), "%v", diags)
	})
}

func TestKMSSecretResourceDelete_pollError(t *testing.T) {
	r := kms.NewSecretResource()
	var schema fwresource.SchemaResponse
	r.Schema(context.Background(), fwresource.SchemaRequest{}, &schema)
	d := newTestState(t, tfsdk.State{Schema: schema.Schema}, map[string]interface{}{kms.SecretFieldDeleteProtection: true, kms.SecretFieldPath: "test-secret"})
	require.False(t, d.SetAttribute(context.Background(), path.Root(kms.EntityFieldID), "test-secret").HasError())
	config := newKMSResourceTestConfig(t,
		kmsResourceRequest{"DELETE", "secret/metadata/test-secret", "", 204, ""},
		kmsResourceRequest{"GET", "secret/data/test-secret", "", 403, `{"errors":["permission denied"]}`},
	)
	r.(fwresource.ResourceWithConfigure).Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: config}, &fwresource.ConfigureResponse{})
	response := fwresource.DeleteResponse{State: d}
	r.Delete(context.Background(), fwresource.DeleteRequest{State: d}, &response)
	d = response.State
	diags := response.Diagnostics
	require.True(t, diags.HasError())
	assert.Equal(t, "Error waiting for secret deletion", diags[0].Summary())
	assert.Contains(t, diags[0].Detail(), "permission denied")
}

func TestAccKMSSecretResource_basic(t *testing.T) {
	name := "tfacc-kms-secret-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccKMSPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		CheckDestroy:             testAccKMSCheckDestroy,
		Steps: []resource.TestStep{
			{Config: fmt.Sprintf(testAccKMSSecretResourceBasic, name), Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "id", name),
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "data_json", `{"password":"initial"}`),
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "delete_protection", "true"),
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "version", "1"),
				resource.TestCheckResourceAttrSet("vkcs_kms_secret.secret", "created_time"),
			)},
			{Config: fmt.Sprintf(testAccKMSSecretResourceConfig, name, "updated", false), Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "id", name),
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "data_json", `{"password":"updated"}`),
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "delete_protection", "false"),
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "version", "2"),
			)},
			{Config: fmt.Sprintf(testAccKMSSecretResourceConfig, name+"-replacement", "replacement", false), Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "id", name+"-replacement"),
				resource.TestCheckResourceAttr("vkcs_kms_secret.secret", "version", "1"),
			)},
		},
	})
}

const testAccKMSSecretResourceBasic = `
resource "vkcs_kms_secret" "secret" {
  path      = %q
  data_json = jsonencode({ password = "initial" })
}
`

const testAccKMSSecretResourceConfig = `
resource "vkcs_kms_secret" "secret" {
  path              = %q
  data_json         = jsonencode({ password = %q })
  delete_protection = %t
}
`
