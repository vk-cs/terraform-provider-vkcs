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
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/planmodifiers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dc/v2/vrrps"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

// Ensure the implementation satisfies the desired interfaces.
var (
	_ resource.Resource                = &VRRPResource{}
	_ resource.ResourceWithConfigure   = &VRRPResource{}
	_ resource.ResourceWithImportState = &VRRPResource{}
)

func NewVRRPResource() resource.Resource {
	return &VRRPResource{}
}

type VRRPResource struct {
	config clients.Config
}

func (r *VRRPResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_dc_vrrp"
}

type VRRPResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	GroupID        types.Int64  `tfsdk:"group_id"`
	NetworkID      types.String `tfsdk:"network_id"`
	SubnetID       types.String `tfsdk:"subnet_id"`
	AdvertInterval types.Int64  `tfsdk:"advert_interval"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	SDN            types.String `tfsdk:"sdn"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`

	Region types.String `tfsdk:"region"`
}

func (r *VRRPResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource",
			},

			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the VRRP",
			},

			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the VRRP",
			},

			"group_id": schema.Int64Attribute{
				Required:    true,
				Description: "VRRP Group ID",
			},

			"network_id": schema.StringAttribute{
				Required:    true,
				Description: "Network ID to attach. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"subnet_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Subnet ID to attach. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},

			"advert_interval": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "VRRP Advertise interval. Default is 1",
			},

			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable or disable item. Default is true",
			},

			"sdn": schema.StringAttribute{
				Computed:    true,
				Description: "SDN of created VRRP",
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
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(planmodifiers.GetRegionPlanModifier(resp),
						"require replacement if configuration value changes", "require replacement if configuration value changes"),
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`. Changing this creates a new vrrp.",
			},
		},
		Description: "Manages a direct connect VRRP resource._note_ This resource requires Sprut SDN to be enabled in your project.",
	}
}

func (r *VRRPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *VRRPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VRRPResourceModel
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

	vrrpCreateOpts := vrrps.CreateOpts{
		Name:           data.Name.ValueString(),
		Description:    data.Description.ValueString(),
		GroupID:        int(data.GroupID.ValueInt64()),
		NetworkID:      data.NetworkID.ValueString(),
		SubnetID:       data.SubnetID.ValueString(),
		AdvertInterval: int(data.AdvertInterval.ValueInt64()),
		Enabled:        util.ValueKnownBoolPointer(data.Enabled),
	}

	vrrpResp, err := vrrps.Create(networkingClient, &vrrps.VRRPCreate{VRRP: &vrrpCreateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_dc_vrrp", err.Error())
		return
	}
	vrrpID := vrrpResp.ID
	resp.State.SetAttribute(ctx, path.Root("id"), vrrpID)

	data.ID = types.StringValue(vrrpResp.ID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(vrrpResp.Name)
	data.Description = types.StringValue(vrrpResp.Description)
	data.GroupID = types.Int64Value(int64(vrrpResp.GroupID))
	data.NetworkID = types.StringValue(vrrpResp.NetworkID)
	data.SubnetID = types.StringValue(vrrpResp.SubnetID)
	data.AdvertInterval = types.Int64Value(int64(vrrpResp.AdvertInterval))
	data.Enabled = types.BoolValue(vrrpResp.Enabled)
	data.SDN = types.StringValue(vrrpResp.SDN)
	data.CreatedAt = types.StringValue(vrrpResp.CreatedAt)
	data.UpdatedAt = types.StringValue(vrrpResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *VRRPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VRRPResourceModel

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

	vrrpID := data.ID.ValueString()

	vrrpResp, err := vrrps.Get(networkingClient, vrrpID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_dc_vrrp", checkDeleted.Error())
		}
		return
	}

	data.ID = types.StringValue(vrrpResp.ID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(vrrpResp.Name)
	data.Description = types.StringValue(vrrpResp.Description)
	data.GroupID = types.Int64Value(int64(vrrpResp.GroupID))
	data.NetworkID = types.StringValue(vrrpResp.NetworkID)
	data.SubnetID = types.StringValue(vrrpResp.SubnetID)
	data.AdvertInterval = types.Int64Value(int64(vrrpResp.AdvertInterval))
	data.Enabled = types.BoolValue(vrrpResp.Enabled)
	data.SDN = types.StringValue(vrrpResp.SDN)
	data.CreatedAt = types.StringValue(vrrpResp.CreatedAt)
	data.UpdatedAt = types.StringValue(vrrpResp.UpdatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VRRPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VRRPResourceModel
	var data VRRPResourceModel

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

	vrrpID := data.ID.ValueString()

	vrrpUpdateOpts := vrrps.UpdateOpts{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		GroupID:        int(plan.GroupID.ValueInt64()),
		AdvertInterval: int(plan.AdvertInterval.ValueInt64()),
		Enabled:        util.ValueKnownBoolPointer(plan.Enabled),
	}

	vrrpResp, err := vrrps.Update(networkingClient, vrrpID, &vrrps.VRRPUpdate{VRRP: &vrrpUpdateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_dc_vrrp", err.Error())
		return
	}

	data.ID = types.StringValue(vrrpResp.ID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(vrrpResp.Name)
	data.Description = types.StringValue(vrrpResp.Description)
	data.GroupID = types.Int64Value(int64(vrrpResp.GroupID))
	data.NetworkID = types.StringValue(vrrpResp.NetworkID)
	data.SubnetID = types.StringValue(vrrpResp.SubnetID)
	data.AdvertInterval = types.Int64Value(int64(vrrpResp.AdvertInterval))
	data.Enabled = types.BoolValue(vrrpResp.Enabled)
	data.SDN = types.StringValue(vrrpResp.SDN)
	data.CreatedAt = types.StringValue(vrrpResp.CreatedAt)
	data.UpdatedAt = types.StringValue(vrrpResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *VRRPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VRRPResourceModel

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

	err = vrrps.Delete(networkingClient, id).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_dc_vrrp", err.Error())
		return
	}
}

func (r *VRRPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
