package dc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dc/v2/bgpstaticannounces"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

// Ensure the implementation satisfies the desired interfaces.
var _ resource.Resource = &BGPStaticAnnounceResource{}

func NewBGPStaticAnnounceResource() resource.Resource {
	return &BGPStaticAnnounceResource{}
}

type BGPStaticAnnounceResource struct {
	config clients.Config
}

func (r *BGPStaticAnnounceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_dc_bgp_static_announce"
}

type BGPStaticAnnounceResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	DCBGPID     types.String `tfsdk:"dc_bgp_id"`
	Network     types.String `tfsdk:"network"`
	Gateway     types.String `tfsdk:"gateway"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	Region      types.String `tfsdk:"region"`
}

func (r *BGPStaticAnnounceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource",
			},

			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the BGP neighbor",
			},

			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the BGP neighbor",
			},

			"dc_bgp_id": schema.StringAttribute{
				Required:    true,
				Description: "Direct Connect BGP ID to attach. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"network": schema.StringAttribute{
				Required:    true,
				Description: "Subnet in CIDR notation. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"gateway": schema.StringAttribute{
				Required:    true,
				Description: "IP address of gateway. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable or disable item. Default is true",
			},

			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Creation timestamp",
			},

			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Update timestamp",
			},

			"region": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`.",
			},
		},
		Description: "Manages a direct connect BGP Static Announce resource.<br> ~> **Note:** This resource requires Sprut SDN to be enabled in your project.",
	}
}

func (r *BGPStaticAnnounceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *BGPStaticAnnounceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BGPStaticAnnounceResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	networkingClient, err := r.config.NetworkingV2Client(region, networking.SprutSDN)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS networking client", err.Error())
		return
	}

	bgpStaticAnnounceCreateOpts := bgpstaticannounces.CreateOpts{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		DCBGPID:     data.DCBGPID.ValueString(),
		Network:     data.Network.ValueString(),
		Gateway:     data.Gateway.ValueString(),
		Enabled:     util.ValueKnownBoolPointer(data.Enabled),
	}

	bgpStaticAnnounceResp, err := bgpstaticannounces.Create(networkingClient, &bgpstaticannounces.BGPStaticAnnounceCreate{BGPStaticAnnounce: &bgpStaticAnnounceCreateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_dc_bgp_static_announce", err.Error())
		return
	}
	bgpStaticAnnounceID := bgpStaticAnnounceResp.ID
	resp.State.SetAttribute(ctx, path.Root("id"), bgpStaticAnnounceID)

	data.ID = types.StringValue(bgpStaticAnnounceID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(bgpStaticAnnounceResp.Name)
	data.Description = types.StringValue(bgpStaticAnnounceResp.Description)
	data.DCBGPID = types.StringValue(bgpStaticAnnounceResp.DCBGPID)
	data.Network = types.StringValue(bgpStaticAnnounceResp.Network)
	data.Gateway = types.StringValue(bgpStaticAnnounceResp.Gateway)
	data.Enabled = types.BoolValue(bgpStaticAnnounceResp.Enabled)
	data.CreatedAt = types.StringValue(bgpStaticAnnounceResp.CreatedAt)
	data.UpdatedAt = types.StringValue(bgpStaticAnnounceResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *BGPStaticAnnounceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BGPStaticAnnounceResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	networkingClient, err := r.config.NetworkingV2Client(region, networking.SprutSDN)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS networking client", err.Error())
		return
	}

	bgpStaticAnnounceID := data.ID.ValueString()

	bgpStaticAnnounceResp, err := bgpstaticannounces.Get(networkingClient, bgpStaticAnnounceID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_dc_bgp_static_announce", checkDeleted.Error())
		}
		return
	}

	data.ID = types.StringValue(bgpStaticAnnounceID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(bgpStaticAnnounceResp.Name)
	data.Description = types.StringValue(bgpStaticAnnounceResp.Description)
	data.DCBGPID = types.StringValue(bgpStaticAnnounceResp.DCBGPID)
	data.Network = types.StringValue(bgpStaticAnnounceResp.Network)
	data.Gateway = types.StringValue(bgpStaticAnnounceResp.Gateway)
	data.Enabled = types.BoolValue(bgpStaticAnnounceResp.Enabled)
	data.CreatedAt = types.StringValue(bgpStaticAnnounceResp.CreatedAt)
	data.UpdatedAt = types.StringValue(bgpStaticAnnounceResp.UpdatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BGPStaticAnnounceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BGPStaticAnnounceResourceModel
	var data BGPStaticAnnounceResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := plan.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	networkingClient, err := r.config.NetworkingV2Client(region, networking.SprutSDN)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS networking client", err.Error())
		return
	}

	bgpStaticAnnounceID := data.ID.ValueString()

	bgpStaticAnnounceUpdateOpts := bgpstaticannounces.UpdateOpts{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Enabled:     util.ValueKnownBoolPointer(plan.Enabled),
	}

	bgpStaticAnnounceResp, err := bgpstaticannounces.Update(networkingClient, bgpStaticAnnounceID, &bgpstaticannounces.BGPStaticAnnounceUpdate{BGPStaticAnnounce: &bgpStaticAnnounceUpdateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_dc_bgp_static_announce", err.Error())
		return
	}

	data.ID = types.StringValue(bgpStaticAnnounceID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(bgpStaticAnnounceResp.Name)
	data.Description = types.StringValue(bgpStaticAnnounceResp.Description)
	data.DCBGPID = types.StringValue(bgpStaticAnnounceResp.DCBGPID)
	data.Network = types.StringValue(bgpStaticAnnounceResp.Network)
	data.Gateway = types.StringValue(bgpStaticAnnounceResp.Gateway)
	data.Enabled = types.BoolValue(bgpStaticAnnounceResp.Enabled)
	data.CreatedAt = types.StringValue(bgpStaticAnnounceResp.CreatedAt)
	data.UpdatedAt = types.StringValue(bgpStaticAnnounceResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *BGPStaticAnnounceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BGPStaticAnnounceResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	networkingClient, err := r.config.NetworkingV2Client(region, networking.SprutSDN)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS networking client", err.Error())
		return
	}

	id := data.ID.ValueString()

	err = bgpstaticannounces.Delete(networkingClient, id).ExtractErr()
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_dc_bgp_static_announce", err.Error())
		return
	}
}

func (r *BGPStaticAnnounceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
