package kms

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/secrets"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
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
			"id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Identifier of the KMS secret"},
			"path": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Description:   "Identifier (path) of the secret within the KMS secret store. Changing this forces a new resource.",
			},
			"data_json": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Data stored in the secret as a JSON object, represented as a string. Use `jsonencode` to encode a Terraform object.",
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time returned in the secret metadata, formatted as a Go time string, for example `2026-09-01 10:00:00 +0000 UTC`.",
			},
			"version": schema.Int64Attribute{
				Computed:    true,
				Description: "Current version of the secret returned in its metadata.",
			},
			"delete_protection": schema.BoolAttribute{
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
		resp.Diagnostics.AddError("Error creating VKCS KMS client", err.Error())
		return
	}

	path := data.Path.ValueString()

	dataJSON := data.DataJSON.ValueString()

	secretData := make(map[string]any)
	err = json.Unmarshal([]byte(dataJSON), &secretData)
	if err != nil {
		resp.Diagnostics.AddError("Error unmarshalling data", err.Error())
		return
	}

	secret, err := createOrUpdateSecret(kmsV1Client, path, secrets.CreateOrUpdateOpts{
		Data: secretData,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating secret", err.Error())
		return
	}

	deleteProtection := data.DeleteProtection.ValueBool()

	// We need to wait for the secret to be ready, then set this flag
	err = waitForKMSOperation(ctx, func() (bool, error) {
		err := setSecretDeleteProtection(kmsV1Client, path, secrets.SetDeleteProtectionOpts{
			DeleteProtection: deleteProtection,
		})
		return err == nil, err
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating secret delete protection", err.Error())
		return
	}

	data.ID = types.StringValue(path)

	_, err = kmsStateChangeConf(secretDeleteProtectionStateRefreshFunc(kmsV1Client, path, &deleteProtection), kmsStateReady).WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for secret delete protection", err.Error())
		return
	}

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
		resp.Diagnostics.AddError("Error creating VKCS KMS client", err.Error())
		return
	}

	path := data.Path.ValueString()

	result, err := kmsStateChangeConf(secretStateRefreshFunc(kmsV1Client, path), kmsStateReady).WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error getting secret", err.Error())
		return
	}

	secret, ok := result.(secrets.GetSecret)
	if !ok {
		resp.Diagnostics.AddError("Unexpected KMS response type", fmt.Sprintf("Expected secrets.GetSecret, got %T", result))
		return
	}

	secretString, err := json.Marshal(secret.Data.Data)
	if err != nil {
		resp.Diagnostics.AddError("Error creating data_json", err.Error())
		return
	}

	protectionResult, err := kmsStateChangeConf(secretDeleteProtectionStateRefreshFunc(kmsV1Client, path, nil), kmsStateReady).WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error getting secret delete protection", err.Error())
		return
	}
	deleteProtection, ok := protectionResult.(bool)
	if !ok {
		resp.Diagnostics.AddError("Unexpected KMS response type", fmt.Sprintf("Expected bool, got %T", protectionResult))
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
		resp.Diagnostics.AddError("Error creating VKCS KMS client", err.Error())
		return
	}

	path := data.Path.ValueString()

	dataJSON := data.DataJSON.ValueString()

	secretData := make(map[string]any)
	err = json.Unmarshal([]byte(dataJSON), &secretData)
	if err != nil {
		resp.Diagnostics.AddError("Error unmarshalling data", err.Error())
		return
	}

	secret, err := createOrUpdateSecret(kmsV1Client, path, secrets.CreateOrUpdateOpts{
		Data: secretData,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating secret", err.Error())
		return
	}

	deleteProtection := data.DeleteProtection.ValueBool()

	// We need to wait for the secret to be ready, then set this flag
	err = waitForKMSOperation(ctx, func() (bool, error) {
		err := setSecretDeleteProtection(kmsV1Client, path, secrets.SetDeleteProtectionOpts{
			DeleteProtection: deleteProtection,
		})
		return err == nil, err
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating secret delete protection", err.Error())
		return
	}

	_, err = kmsStateChangeConf(secretDeleteProtectionStateRefreshFunc(kmsV1Client, path, &deleteProtection), kmsStateReady).WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for secret delete protection", err.Error())
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
		resp.Diagnostics.AddError("Error creating VKCS KMS client", err.Error())
		return
	}

	path := data.Path.ValueString()

	err = waitForKMSOperation(ctx, func() (bool, error) {
		err := deleteSecret(kmsV1Client, path)
		if errutil.IsNotFound(err) {
			return true, nil
		}
		return err == nil, err
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting secret", err.Error())
		return
	}

	_, err = kmsStateChangeConf(secretStateRefreshFunc(kmsV1Client, path), kmsStateMissing).WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for secret deletion", err.Error())
	}
}
