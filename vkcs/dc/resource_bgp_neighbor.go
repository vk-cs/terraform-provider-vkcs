package dc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dc/v2/bgpneighbors"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

// Ensure the implementation satisfies the desired interfaces.
var _ resource.Resource = &BGPNeighborResource{}

func NewBGPNeighborResource() resource.Resource {
	return &BGPNeighborResource{}
}

type BGPNeighborResource struct {
	config clients.Config
}

func (r *BGPNeighborResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_dc_bgp_neighbor"
}

type BGPNeighborResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	DCBGPID              types.String `tfsdk:"dc_bgp_id"`
	RemoteASN            types.Int64  `tfsdk:"remote_asn"`
	RemoteIP             types.String `tfsdk:"remote_ip"`
	ForceIBGPNextHopSelf types.Bool   `tfsdk:"force_ibgp_next_hop_self"`
	AddPaths             types.String `tfsdk:"add_paths"`
	BFDEnabled           types.Bool   `tfsdk:"bfd_enabled"`
	FilterIn             types.String `tfsdk:"filter_in"`
	FilterOut            types.String `tfsdk:"filter_out"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
	Region               types.String `tfsdk:"region"`
}

func (r *BGPNeighborResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"remote_asn": schema.Int64Attribute{
				Required:    true,
				Description: "BGP Neighbor ASN. Changing this creates a new resource",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},

			"remote_ip": schema.StringAttribute{
				Required:    true,
				Description: "BGP Neighbor IP address. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"force_ibgp_next_hop_self": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Force set IP address of next-hop on BGP prefix to self even in iBGP. Default is false",
			},

			"add_paths": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Activate BGP Add-Paths feature on peer. Default is off",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"off", "on", "tx", "rx"}...),
				},
			},

			"bfd_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Control BGP session activity with BFD protocol. Default is false",
			},

			"filter_in": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Input filter that pass incoming BGP prefixes (allow any)",
			},

			"filter_out": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Output filter that pass incoming BGP prefixes (allow any)",
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
		Description: "Manages a direct connect BGP neighbor resource. **Note:** This resource requires Sprut SDN to be enabled in your project.",
	}
}

func (r *BGPNeighborResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *BGPNeighborResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BGPNeighborResourceModel
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

	bgpNeighborCreateOpts := bgpneighbors.CreateOpts{
		Name:                 data.Name.ValueString(),
		Description:          data.Description.ValueString(),
		DCBGPID:              data.DCBGPID.ValueString(),
		RemoteASN:            int(data.RemoteASN.ValueInt64()),
		RemoteIP:             data.RemoteIP.ValueString(),
		AddPaths:             data.AddPaths.ValueString(),
		FilterIn:             data.FilterIn.ValueString(),
		FilterOut:            data.FilterOut.ValueString(),
		ForceIBGPNextHopSelf: util.ValueKnownBoolPointer(data.ForceIBGPNextHopSelf),
		BFDEnabled:           util.ValueKnownBoolPointer(data.BFDEnabled),
		Enabled:              util.ValueKnownBoolPointer(data.Enabled),
	}

	bgpNeighborResp, err := bgpneighbors.Create(networkingClient, &bgpneighbors.BGPNeighborCreate{BGPNeighbor: &bgpNeighborCreateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_dc_bgp_neighbor", err.Error())
		return
	}
	bgpNeighborID := bgpNeighborResp.ID
	resp.State.SetAttribute(ctx, path.Root("id"), bgpNeighborID)

	data.ID = types.StringValue(bgpNeighborID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(bgpNeighborResp.Name)
	data.Description = types.StringValue(bgpNeighborResp.Description)
	data.DCBGPID = types.StringValue(bgpNeighborResp.DCBGPID)
	data.RemoteASN = types.Int64Value(int64(bgpNeighborResp.RemoteASN))
	data.RemoteIP = types.StringValue(bgpNeighborResp.RemoteIP)
	data.ForceIBGPNextHopSelf = types.BoolValue(bgpNeighborResp.ForceIBGPNextHopSelf)
	data.AddPaths = types.StringValue(bgpNeighborResp.AddPaths)
	data.BFDEnabled = types.BoolValue(bgpNeighborResp.BFDEnabled)
	data.FilterIn = types.StringValue(bgpNeighborResp.FilterIn)
	data.FilterOut = types.StringValue(bgpNeighborResp.FilterOut)
	data.Enabled = types.BoolValue(bgpNeighborResp.Enabled)
	data.CreatedAt = types.StringValue(bgpNeighborResp.CreatedAt)
	data.UpdatedAt = types.StringValue(bgpNeighborResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *BGPNeighborResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BGPNeighborResourceModel

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

	bgpNeighborID := data.ID.ValueString()

	bgpNeighborResp, err := bgpneighbors.Get(networkingClient, bgpNeighborID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_dc_bgp_neighbor", checkDeleted.Error())
		}
		return
	}

	data.ID = types.StringValue(bgpNeighborID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(bgpNeighborResp.Name)
	data.Description = types.StringValue(bgpNeighborResp.Description)
	data.DCBGPID = types.StringValue(bgpNeighborResp.DCBGPID)
	data.RemoteASN = types.Int64Value(int64(bgpNeighborResp.RemoteASN))
	data.RemoteIP = types.StringValue(bgpNeighborResp.RemoteIP)
	data.ForceIBGPNextHopSelf = types.BoolValue(bgpNeighborResp.ForceIBGPNextHopSelf)
	data.AddPaths = types.StringValue(bgpNeighborResp.AddPaths)
	data.BFDEnabled = types.BoolValue(bgpNeighborResp.BFDEnabled)
	data.FilterIn = types.StringValue(bgpNeighborResp.FilterIn)
	data.FilterOut = types.StringValue(bgpNeighborResp.FilterOut)
	data.Enabled = types.BoolValue(bgpNeighborResp.Enabled)
	data.CreatedAt = types.StringValue(bgpNeighborResp.CreatedAt)
	data.UpdatedAt = types.StringValue(bgpNeighborResp.UpdatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BGPNeighborResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BGPNeighborResourceModel
	var data BGPNeighborResourceModel

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

	bgpNeighborID := data.ID.ValueString()

	bgpNeighborUpdateOpts := bgpneighbors.UpdateOpts{
		Name:                 plan.Name.ValueString(),
		Description:          plan.Description.ValueString(),
		AddPaths:             plan.AddPaths.ValueString(),
		FilterIn:             plan.FilterIn.ValueString(),
		FilterOut:            plan.FilterOut.ValueString(),
		ForceIBGPNextHopSelf: util.ValueKnownBoolPointer(plan.ForceIBGPNextHopSelf),
		BFDEnabled:           util.ValueKnownBoolPointer(plan.BFDEnabled),
		Enabled:              util.ValueKnownBoolPointer(plan.Enabled),
	}

	bgpNeighborResp, err := bgpneighbors.Update(networkingClient, bgpNeighborID, &bgpneighbors.BGPNeighborUpdate{BGPNeighbor: &bgpNeighborUpdateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_dc_bgp_instance", err.Error())
		return
	}

	data.ID = types.StringValue(bgpNeighborID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(bgpNeighborResp.Name)
	data.Description = types.StringValue(bgpNeighborResp.Description)
	data.DCBGPID = types.StringValue(bgpNeighborResp.DCBGPID)
	data.RemoteASN = types.Int64Value(int64(bgpNeighborResp.RemoteASN))
	data.RemoteIP = types.StringValue(bgpNeighborResp.RemoteIP)
	data.ForceIBGPNextHopSelf = types.BoolValue(bgpNeighborResp.ForceIBGPNextHopSelf)
	data.AddPaths = types.StringValue(bgpNeighborResp.AddPaths)
	data.BFDEnabled = types.BoolValue(bgpNeighborResp.BFDEnabled)
	data.FilterIn = types.StringValue(bgpNeighborResp.FilterIn)
	data.FilterOut = types.StringValue(bgpNeighborResp.FilterOut)
	data.Enabled = types.BoolValue(bgpNeighborResp.Enabled)
	data.CreatedAt = types.StringValue(bgpNeighborResp.CreatedAt)
	data.UpdatedAt = types.StringValue(bgpNeighborResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *BGPNeighborResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BGPNeighborResourceModel

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

	err = bgpneighbors.Delete(networkingClient, id).ExtractErr()
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_dc_bgp_neighbor", err.Error())
		return
	}
}

func (r *BGPNeighborResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
