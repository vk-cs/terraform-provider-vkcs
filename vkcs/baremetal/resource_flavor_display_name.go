package baremetal

import (
	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/baremetal/v1/flavors"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

var (
	_ resource.Resource                = (*FlavorDisplayNameResource)(nil)
	_ resource.ResourceWithConfigure   = (*FlavorDisplayNameResource)(nil)
	_ resource.ResourceWithImportState = (*FlavorDisplayNameResource)(nil)
)

// NewFlavorDisplayNameResource manages the project-specific display name of an existing bare metal flavor.
func NewFlavorDisplayNameResource() resource.Resource {
	return &FlavorDisplayNameResource{}
}

type FlavorDisplayNameResource struct {
	config clients.Config
}

type FlavorDisplayNameResourceModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
}

func (r *FlavorDisplayNameResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_baremetal_flavor_display_name"
}

func (r *FlavorDisplayNameResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The UUID of the existing flavor.",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The project-specific display name of the flavor.",
			},
		},
		Description: "Manages the project-specific display name of an existing VKCS bare metal flavor.",
	}
}

func (r *FlavorDisplayNameResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.config = req.ProviderData.(clients.Config)
}

func (r *FlavorDisplayNameResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FlavorDisplayNameResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.config.BareMetalV1Client(r.config.GetRegion())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS baremetal API client", err.Error())
		return
	}

	if err := setFlavorDisplayName(client, data.ID.ValueString(), data.DisplayName.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error updating baremetal flavor display name", err.Error())
		return
	}

	flavor, err := flavors.Get(client, data.ID.ValueString()).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error reading baremetal flavor after update", err.Error())
		return
	}

	setFlavorDisplayNameResourceModel(&data, flavor)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FlavorDisplayNameResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FlavorDisplayNameResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.config.BareMetalV1Client(r.config.GetRegion())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS baremetal API client", err.Error())
		return
	}

	flavor, err := flavors.Get(client, data.ID.ValueString()).Extract()
	if errutil.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading baremetal flavor", err.Error())
		return
	}

	setFlavorDisplayNameResourceModel(&data, flavor)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FlavorDisplayNameResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FlavorDisplayNameResourceModel
	var state FlavorDisplayNameResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.config.BareMetalV1Client(r.config.GetRegion())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS baremetal API client", err.Error())
		return
	}

	if err := setFlavorDisplayName(client, state.ID.ValueString(), plan.DisplayName.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error updating baremetal flavor display name", err.Error())
		return
	}

	flavor, err := flavors.Get(client, state.ID.ValueString()).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error reading baremetal flavor after update", err.Error())
		return
	}

	setFlavorDisplayNameResourceModel(&state, flavor)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FlavorDisplayNameResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FlavorDisplayNameResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.config.BareMetalV1Client(r.config.GetRegion())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS baremetal API client", err.Error())
		return
	}

	if err := setFlavorDisplayName(client, data.ID.ValueString(), ""); errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error clearing baremetal flavor display name", err.Error())
	}
}

func (r *FlavorDisplayNameResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func setFlavorDisplayName(client *gophercloud.ServiceClient, flavorID, displayName string) error {
	return flavors.Update(client, flavorID, flavors.UpdateOpts{DisplayName: &displayName}).ExtractErr()
}

func setFlavorDisplayNameResourceModel(data *FlavorDisplayNameResourceModel, flavor *flavors.Flavor) {
	data.ID = types.StringValue(flavor.FlavorId)
	data.DisplayName = types.StringPointerValue(flavor.DisplayName)
}
