package kms

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/secrets"
)

var _ resource.ResourceWithConfigure = (*secretResource)(nil)

func NewSecretResource() resource.Resource {
	return &secretResource{}
}

type secretResource struct {
	config clients.Config
}

type secretResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Path             types.String `tfsdk:"path"`
	DataJSON         types.String `tfsdk:"data_json"`
	CreatedTime      types.String `tfsdk:"created_time"`
	Version          types.Int64  `tfsdk:"version"`
	DeleteProtection types.Bool   `tfsdk:"delete_protection"`
}

func (r *secretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kms_secret"
}

func (r *secretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		r.config = req.ProviderData.(clients.Config)
	}
}

func (r *secretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource representing KMS secret",
		Attributes: map[string]schema.Attribute{
			EntityFieldID: schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Identifier of the KMS secret"},
			SecretFieldPath: schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Description:   "Identifier (path) of the secret within the KMS secret store. Changing this forces a new resource.",
			},
			SecretFieldDataJSON: schema.StringAttribute{
				Required:    true,
				Description: "Data stored in the secret as a JSON object, represented as a string. Use `jsonencode` to encode a Terraform object. Refreshed from KMS when the resource is read.",
			},
			SecretFieldCreatedTime: schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Creation time returned in the secret metadata, formatted as a Go time string, for example `2026-09-01 10:00:00 +0000 UTC`.",
			},
			SecretFieldVersion: schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Current version of the secret returned in its metadata.",
			},
			SecretFieldDeleteProtection: schema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the secret is protected from deletion. Defaults to `true`. Set to `false` and apply before destroying the secret.",
			},
		},
	}
}

func (r *secretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data secretResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := r.config
	kmsV1Client, err := config.KMSV1Client(config.GetRegion())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error creating VKCS KMS client: %s", err), "")
		return
	}

	path := data.Path.ValueString()

	dataJSON := data.DataJSON.ValueString()

	secretData := make(map[string]any)
	err = json.Unmarshal([]byte(dataJSON), &secretData)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error unmarshalling data: %s", err), "")
		return
	}

	secret, err := createOrUpdateSecret(kmsV1Client, path, secrets.CreateOrUpdateOpts{
		Data: secretData,
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error creating secret: %s", err), "")
		return
	}

	deleteProtection := data.DeleteProtection.ValueBool()

	// We need to wait for the secret to be ready, then set this flag
	for range retriesCount {
		err = setSecretDeleteProtection(kmsV1Client, path, secrets.SetDeleteProtectionOpts{
			DeleteProtection: deleteProtection,
		})
		if err == nil {
			break
		}
		time.Sleep(defaultThreshold)
	}
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error updating secret delete protection: %s", err), "")
		return
	}

	data.ID = types.StringValue(path)

	data.CreatedTime = types.StringValue(secret.Data.CreatedTime.String())

	data.Version = types.Int64Value(int64(secret.Data.Version))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *secretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data secretResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := r.config
	kmsV1Client, err := config.KMSV1Client(config.GetRegion())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error creating VKCS KMS client: %s", err), "")
		return
	}

	path := data.Path.ValueString()

	secret, err := getSecret(kmsV1Client, path)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error getting secret: %s", err), "")
		return
	}

	secretString, err := json.Marshal(secret.Data.Data)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error creating data_json: %s", err), "")
		return
	}

	deleteProtection, err := getSecretDeleteProtection(kmsV1Client, path)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error getting secret delete protection: %s", err), "")
		return
	}

	data.DataJSON = types.StringValue(string(secretString))

	data.CreatedTime = types.StringValue(secret.Data.Metadata.CreatedTime.String())

	data.Version = types.Int64Value(int64(secret.Data.Metadata.Version))

	data.DeleteProtection = types.BoolValue(deleteProtection)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *secretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data secretResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := r.config
	kmsV1Client, err := config.KMSV1Client(config.GetRegion())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error creating VKCS KMS client: %s", err), "")
		return
	}

	path := data.Path.ValueString()

	dataJSON := data.DataJSON.ValueString()

	secretData := make(map[string]any)
	err = json.Unmarshal([]byte(dataJSON), &secretData)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error unmarshalling data: %s", err), "")
		return
	}

	secret, err := createOrUpdateSecret(kmsV1Client, path, secrets.CreateOrUpdateOpts{
		Data: secretData,
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error creating secret: %s", err), "")
		return
	}

	deleteProtection := data.DeleteProtection.ValueBool()

	// We need to wait for the secret to be ready, then set this flag
	for range retriesCount {
		err = setSecretDeleteProtection(kmsV1Client, path, secrets.SetDeleteProtectionOpts{
			DeleteProtection: deleteProtection,
		})
		if err == nil {
			break
		}
		time.Sleep(defaultThreshold)
	}
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error updating secret delete protection: %s", err), "")
		return
	}

	data.CreatedTime = types.StringValue(secret.Data.CreatedTime.String())

	data.Version = types.Int64Value(int64(secret.Data.Version))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *secretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data secretResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := r.config
	kmsV1Client, err := config.KMSV1Client(config.GetRegion())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error creating VKCS KMS client: %s", err), "")
		return
	}

	path := data.Path.ValueString()

	err = deleteSecret(kmsV1Client, path)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error deleting secret: %s", err), "")
		return
	}

	// Waiting for the secret to be fully deleted
	for range retriesCount {
		var err error
		_, err = getSecret(kmsV1Client, path)
		if err != nil {
			break
		}
		time.Sleep(defaultThreshold)
	}
}
