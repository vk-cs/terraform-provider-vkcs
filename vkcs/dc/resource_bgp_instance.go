package dc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	bgpinstances "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dc/v2/bgpinstances"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

// Ensure the implementation satisfies the desired interfaces.
var _ resource.Resource = &BGPInstanceResource{}

func NewBGPInstanceResource() resource.Resource {
	return &BGPInstanceResource{}
}

type BGPInstanceResource struct {
	config clients.Config
}

func (r *BGPInstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_dc_bgp_instance"
}

type BGPInstanceResourceModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Description              types.String `tfsdk:"description"`
	DCRouterID               types.String `tfsdk:"dc_router_id"`
	BGPRouterID              types.String `tfsdk:"bgp_router_id"`
	ASN                      types.Int64  `tfsdk:"asn"`
	ECMPEnabled              types.Bool   `tfsdk:"ecmp_enabled"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	GracefulRestart          types.Bool   `tfsdk:"graceful_restart"`
	LongLivedGracefulRestart types.Bool   `tfsdk:"long_lived_graceful_restart"`
	CreatedAt                types.String `tfsdk:"created_at"`
	UpdatedAt                types.String `tfsdk:"updated_at"`
	Region                   types.String `tfsdk:"region"`
}

func (r *BGPInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource",
			},

			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the router",
			},

			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the router",
			},

			"dc_router_id": schema.StringAttribute{
				Required:    true,
				Description: "Direct Connect Router ID to attach. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"bgp_router_id": schema.StringAttribute{
				Required:    true,
				Description: "BGP Router ID (IP address that represent BGP router in BGP network). Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"asn": schema.Int64Attribute{
				Required:    true,
				Description: "BGP Autonomous System Number (integer representation supports only). Changing this creates a new resource",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},

			"ecmp_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable BGP ECMP behaviour on router. Default is false",
			},

			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable or disable item. Default is true",
			},

			"graceful_restart": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable BGP Graceful Restart feature. Default is false",
			},

			"long_lived_graceful_restart": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable BGP Long Lived Graceful Restart feature. Default is false",
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
		Description: "Manages a direct connect BGP instance resource. **Note:** This resource requires Sprut SDN to be enabled in your project.",
	}
}

func (r *BGPInstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *BGPInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BGPInstanceResourceModel
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

	bgpInstanceCreateOpts := bgpinstances.CreateOpts{
		Name:                     data.Name.ValueString(),
		Description:              data.Description.ValueString(),
		DCRouterID:               data.DCRouterID.ValueString(),
		BGPRouterID:              data.BGPRouterID.ValueString(),
		ASN:                      int(data.ASN.ValueInt64()),
		ECMPEnabled:              util.ValueKnownBoolPointer(data.ECMPEnabled),
		Enabled:                  util.ValueKnownBoolPointer(data.Enabled),
		GracefulRestart:          util.ValueKnownBoolPointer(data.GracefulRestart),
		LongLivedGracefulRestart: util.ValueKnownBoolPointer(data.LongLivedGracefulRestart),
	}

	bgpInstanceResp, err := bgpinstances.Create(networkingClient, &bgpinstances.BGPInstanceCreate{BGPInstance: &bgpInstanceCreateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_dc_bgp_instance", err.Error())
		return
	}
	bgpInstanceID := bgpInstanceResp.ID
	resp.State.SetAttribute(ctx, path.Root("id"), bgpInstanceID)

	data.ID = types.StringValue(bgpInstanceID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(bgpInstanceResp.Name)
	data.Description = types.StringValue(bgpInstanceResp.Description)
	data.DCRouterID = types.StringValue(bgpInstanceResp.DCRouterID)
	data.BGPRouterID = types.StringValue(bgpInstanceResp.BGPRouterID)
	data.ASN = types.Int64Value(int64(bgpInstanceResp.ASN))
	data.ECMPEnabled = types.BoolValue(bgpInstanceResp.ECMPEnabled)
	data.Enabled = types.BoolValue(bgpInstanceResp.Enabled)
	data.GracefulRestart = types.BoolValue(bgpInstanceResp.GracefulRestart)
	data.LongLivedGracefulRestart = types.BoolValue(bgpInstanceResp.LongLivedGracefulRestart)

	data.CreatedAt = types.StringValue(bgpInstanceResp.CreatedAt)
	data.UpdatedAt = types.StringValue(bgpInstanceResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *BGPInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BGPInstanceResourceModel

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

	bgpInstanceID := data.ID.ValueString()

	bgpInstanceResp, err := bgpinstances.Get(networkingClient, bgpInstanceID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_dc_bgp_instance", checkDeleted.Error())
		}
		return
	}

	data.ID = types.StringValue(bgpInstanceID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(bgpInstanceResp.Name)
	data.Description = types.StringValue(bgpInstanceResp.Description)
	data.DCRouterID = types.StringValue(bgpInstanceResp.DCRouterID)
	data.BGPRouterID = types.StringValue(bgpInstanceResp.BGPRouterID)
	data.ASN = types.Int64Value(int64(bgpInstanceResp.ASN))
	data.ECMPEnabled = types.BoolValue(bgpInstanceResp.ECMPEnabled)
	data.Enabled = types.BoolValue(bgpInstanceResp.Enabled)
	data.GracefulRestart = types.BoolValue(bgpInstanceResp.GracefulRestart)
	data.LongLivedGracefulRestart = types.BoolValue(bgpInstanceResp.LongLivedGracefulRestart)
	data.CreatedAt = types.StringValue(bgpInstanceResp.CreatedAt)
	data.UpdatedAt = types.StringValue(bgpInstanceResp.UpdatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BGPInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BGPInstanceResourceModel
	var data BGPInstanceResourceModel

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

	bgpInstanceID := data.ID.ValueString()

	bgpInstanceUpdateOpts := bgpinstances.UpdateOpts{
		Name:                     plan.Name.ValueString(),
		Description:              plan.Description.ValueString(),
		ECMPEnabled:              util.ValueKnownBoolPointer(plan.ECMPEnabled),
		Enabled:                  util.ValueKnownBoolPointer(plan.Enabled),
		GracefulRestart:          util.ValueKnownBoolPointer(plan.GracefulRestart),
		LongLivedGracefulRestart: util.ValueKnownBoolPointer(plan.LongLivedGracefulRestart),
	}

	bgpInstanceResp, err := bgpinstances.Update(networkingClient, bgpInstanceID, &bgpinstances.BGPInstanceUpdate{BGPInstance: &bgpInstanceUpdateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_dc_bgp_instance", err.Error())
		return
	}

	data.ID = types.StringValue(bgpInstanceID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(bgpInstanceResp.Name)
	data.Description = types.StringValue(bgpInstanceResp.Description)
	data.DCRouterID = types.StringValue(bgpInstanceResp.DCRouterID)
	data.BGPRouterID = types.StringValue(bgpInstanceResp.BGPRouterID)
	data.ASN = types.Int64Value(int64(bgpInstanceResp.ASN))
	data.ECMPEnabled = types.BoolValue(bgpInstanceResp.ECMPEnabled)
	data.Enabled = types.BoolValue(bgpInstanceResp.Enabled)
	data.GracefulRestart = types.BoolValue(bgpInstanceResp.GracefulRestart)
	data.LongLivedGracefulRestart = types.BoolValue(bgpInstanceResp.LongLivedGracefulRestart)
	data.CreatedAt = types.StringValue(bgpInstanceResp.CreatedAt)
	data.UpdatedAt = types.StringValue(bgpInstanceResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *BGPInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BGPInstanceResourceModel

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

	err = bgpinstances.Delete(networkingClient, id).ExtractErr()
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_dc_bgp_instance", err.Error())
		return
	}
}

func (r *BGPInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
