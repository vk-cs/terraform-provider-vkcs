package kms

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/keys"
)

var _ resource.ResourceWithConfigure = (*keyResource)(nil)

func NewKeyResource() resource.Resource {
	return &keyResource{}
}

type keyResource struct {
	config clients.Config
}

type keyResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	DeletionAllowed types.Bool   `tfsdk:"deletion_allowed"`
}

func (r *keyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = entityName(req.ProviderTypeName, KeyEntityName)
}

func (r *keyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		r.config = req.ProviderData.(clients.Config)
	}
}

func (r *keyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource representing KMS key",
		Attributes: map[string]schema.Attribute{
			EntityFieldID: schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Identifier of the KMS key"},
			KeyFieldName: schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Description:   "Name of the KMS key. Changing this forces a new resource.",
			},
			KeyFieldType: schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Default:       stringdefault.StaticString("aes256-gcm96"),
				Description:   "Type of the KMS key: `aes128-gcm96`, `aes256-gcm96`, `chacha20-poly1305`, or `xchacha20-poly1305`. Defaults to `aes256-gcm96`. Changing this forces a new resource.",
			},
			KeyFieldDeletionAllowed: schema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether KMS allows the key to be deleted. Defaults to `false`. Set to `true` and apply before destroying the key. Refreshed from KMS when the resource is read.",
			},
		},
	}
}

func (r *keyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data keyResourceModel
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

	name := data.Name.ValueString()

	t := data.Type.ValueString()

	_, err = createKey(kmsV1Client, name, keys.CreateOpts{
		Type: t,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating key", err.Error())
		return
	}

	deletionAllowed := data.DeletionAllowed.ValueBool()

	// We need to wait for the key to be ready, then set this flag
	for range retriesCount {
		_, err = updateKey(kmsV1Client, name, keys.UpdateOpts{
			DeletionAllowed: deletionAllowed,
		})
		if err == nil {
			break
		}
		time.Sleep(defaultThreshold)
	}
	if err != nil {
		resp.Diagnostics.AddError("Error updating key configuration", err.Error())
		return
	}

	data.ID = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data keyResourceModel
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

	name := data.Name.ValueString()

	key, err := getKey(kmsV1Client, name)
	if err != nil {
		resp.Diagnostics.AddError("Error getting key", err.Error())
		return
	}

	data.Name = types.StringValue(key.Data.Name)

	data.Type = types.StringValue(key.Data.Type)

	data.DeletionAllowed = types.BoolValue(key.Data.DeletionAllowed)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data keyResourceModel
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

	deletionAllowed := data.DeletionAllowed.ValueBool()

	_, err = updateKey(kmsV1Client, data.ID.ValueString(), keys.UpdateOpts{
		DeletionAllowed: deletionAllowed,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating key configuration", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data keyResourceModel
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

	name := data.Name.ValueString()

	err = deleteKey(kmsV1Client, name)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting key", err.Error())
		return
	}

	// Waiting for the key to be fully deleted
	for range retriesCount {
		var err error
		_, err = getKey(kmsV1Client, name)
		if err != nil {
			break
		}
		time.Sleep(defaultThreshold)
	}
}
