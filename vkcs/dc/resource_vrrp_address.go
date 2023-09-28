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
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dc/v2/vrrpaddresses"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

// Ensure the implementation satisfies the desired interfaces.
var _ resource.Resource = &VRRPAddressResource{}

func NewVRRPAddressResource() resource.Resource {
	return &VRRPAddressResource{}
}

type VRRPAddressResource struct {
	config clients.Config
}

func (r *VRRPAddressResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_dc_vrrp_address"
}

type VRRPAddressResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	DCVRRPID    types.String `tfsdk:"dc_vrrp_id"`
	IPAddress   types.String `tfsdk:"ip_address"`
	PortID      types.String `tfsdk:"port_id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	Region      types.String `tfsdk:"region"`
}

func (r *VRRPAddressResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"dc_vrrp_id": schema.StringAttribute{
				Required:    true,
				Description: "VRRP ID to attach. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"ip_address": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "IP address to assign. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},

			"port_id": schema.StringAttribute{
				Computed:    true,
				Description: "Port ID used to assign IP address",
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
				Optional:    true,
				Computed:    true,
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`.",
			},
		},
		Description: "Manages a direct connect VRRP address resource.<br> ~> **Note:** This resource requires Sprut SDN to be enabled in your project.",
	}
}

func (r *VRRPAddressResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *VRRPAddressResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VRRPAddressResourceModel
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

	vrrpAddressCreateOpts := vrrpaddresses.CreateOpts{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		DCVRRPID:    data.DCVRRPID.ValueString(),
		IPAddress:   data.IPAddress.ValueString(),
	}

	vrrpAddressResp, err := vrrpaddresses.Create(networkingClient, &vrrpaddresses.VRRPAddressCreate{VRRPAddress: &vrrpAddressCreateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_dc_vrrp_address", err.Error())
		return
	}
	vrrpAddressID := vrrpAddressResp.ID
	resp.State.SetAttribute(ctx, path.Root("id"), vrrpAddressID)

	data.ID = types.StringValue(vrrpAddressResp.ID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(vrrpAddressResp.Name)
	data.Description = types.StringValue(vrrpAddressResp.Description)
	data.DCVRRPID = types.StringValue(vrrpAddressResp.DCVRRPID)
	data.IPAddress = types.StringValue(vrrpAddressResp.IPAddress)
	data.PortID = types.StringValue(vrrpAddressResp.PortID)
	data.CreatedAt = types.StringValue(vrrpAddressResp.CreatedAt)
	data.UpdatedAt = types.StringValue(vrrpAddressResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *VRRPAddressResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VRRPAddressResourceModel

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

	vrrpAddressID := data.ID.ValueString()

	vrrpAddressResp, err := vrrpaddresses.Get(networkingClient, vrrpAddressID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_dc_vrrp_address", checkDeleted.Error())
		}
		return
	}

	data.ID = types.StringValue(vrrpAddressResp.ID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(vrrpAddressResp.Name)
	data.Description = types.StringValue(vrrpAddressResp.Description)
	data.DCVRRPID = types.StringValue(vrrpAddressResp.DCVRRPID)
	data.IPAddress = types.StringValue(vrrpAddressResp.IPAddress)
	data.PortID = types.StringValue(vrrpAddressResp.PortID)
	data.CreatedAt = types.StringValue(vrrpAddressResp.CreatedAt)
	data.UpdatedAt = types.StringValue(vrrpAddressResp.UpdatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VRRPAddressResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VRRPAddressResourceModel
	var data VRRPAddressResourceModel

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

	vrrpAddressID := data.ID.ValueString()

	vrrpAddressUpdateOpts := vrrpaddresses.UpdateOpts{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	vrrpAddressResp, err := vrrpaddresses.Update(networkingClient, vrrpAddressID, &vrrpaddresses.VRRPAddressUpdate{VRRPAddress: &vrrpAddressUpdateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_dc_vrrp_address", err.Error())
		return
	}

	data.ID = types.StringValue(vrrpAddressResp.ID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(vrrpAddressResp.Name)
	data.Description = types.StringValue(vrrpAddressResp.Description)
	data.DCVRRPID = types.StringValue(vrrpAddressResp.DCVRRPID)
	data.IPAddress = types.StringValue(vrrpAddressResp.IPAddress)
	data.PortID = types.StringValue(vrrpAddressResp.PortID)
	data.CreatedAt = types.StringValue(vrrpAddressResp.CreatedAt)
	data.UpdatedAt = types.StringValue(vrrpAddressResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *VRRPAddressResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VRRPAddressResourceModel

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

	err = vrrpaddresses.Delete(networkingClient, id).ExtractErr()
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_dc_vrrp_address", err.Error())
		return
	}
}

func (r *VRRPAddressResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
