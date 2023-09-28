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
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dc/v2/vrrpinterfaces"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

// Ensure the implementation satisfies the desired interfaces.
var _ resource.Resource = &VRRPInterfaceResource{}

func NewVRRPInterfaceResource() resource.Resource {
	return &VRRPInterfaceResource{}
}

type VRRPInterfaceResource struct {
	config clients.Config
}

func (r *VRRPInterfaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_dc_vrrp_interface"
}

type VRRPInterfaceResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	DCVRRPID      types.String `tfsdk:"dc_vrrp_id"`
	DCInterfaceID types.String `tfsdk:"dc_interface_id"`
	Priority      types.Int64  `tfsdk:"priority"`
	Preempt       types.Bool   `tfsdk:"preempt"`
	Master        types.Bool   `tfsdk:"master"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	Region        types.String `tfsdk:"region"`
}

func (r *VRRPInterfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"dc_interface_id": schema.StringAttribute{
				Required:    true,
				Description: "DC Interface ID to attach. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"priority": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "VRRP interface priority. Default is 100",
			},

			"preempt": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "VRRP interface preempt behavior. Default is true",
			},

			"master": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Start VRRP instance on interface as VRRP Master. Default is false",
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
		Description: "Manages a direct connect VRRP interface resource.<br> ~> **Note:** This resource requires Sprut SDN to be enabled in your project.",
	}
}

func (r *VRRPInterfaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *VRRPInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VRRPInterfaceResourceModel
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

	vrrpInterfaceCreateOpts := vrrpinterfaces.CreateOpts{
		Name:          data.Name.ValueString(),
		Description:   data.Description.ValueString(),
		DCVRRPID:      data.DCVRRPID.ValueString(),
		DCInterfaceID: data.DCInterfaceID.ValueString(),
		Priority:      int(data.Priority.ValueInt64()),
		Preempt:       util.ValueKnownBoolPointer(data.Preempt),
		Master:        util.ValueKnownBoolPointer(data.Master),
	}

	vrrpInterfaceResp, err := vrrpinterfaces.Create(networkingClient, &vrrpinterfaces.VRRPInterfaceCreate{VRRPInterface: &vrrpInterfaceCreateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_dc_vrrp_interface", err.Error())
		return
	}
	vrrpInterfaceID := vrrpInterfaceResp.ID
	resp.State.SetAttribute(ctx, path.Root("id"), vrrpInterfaceID)

	data.ID = types.StringValue(vrrpInterfaceResp.ID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(vrrpInterfaceResp.Name)
	data.Description = types.StringValue(vrrpInterfaceResp.Description)
	data.DCVRRPID = types.StringValue(vrrpInterfaceResp.DCVRRPID)
	data.DCInterfaceID = types.StringValue(vrrpInterfaceResp.DCInterfaceID)
	data.Priority = types.Int64Value(int64(vrrpInterfaceResp.Priority))
	data.Preempt = types.BoolValue(vrrpInterfaceResp.Preempt)
	data.Master = types.BoolValue(vrrpInterfaceResp.Master)
	data.CreatedAt = types.StringValue(vrrpInterfaceResp.CreatedAt)
	data.UpdatedAt = types.StringValue(vrrpInterfaceResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *VRRPInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VRRPInterfaceResourceModel

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

	vrrpInterfaceID := data.ID.ValueString()

	vrrpInterfaceResp, err := vrrpinterfaces.Get(networkingClient, vrrpInterfaceID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_dc_vrrp_interface", checkDeleted.Error())
		}
		return
	}

	data.ID = types.StringValue(vrrpInterfaceResp.ID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(vrrpInterfaceResp.Name)
	data.Description = types.StringValue(vrrpInterfaceResp.Description)
	data.DCVRRPID = types.StringValue(vrrpInterfaceResp.DCVRRPID)
	data.DCInterfaceID = types.StringValue(vrrpInterfaceResp.DCInterfaceID)
	data.Priority = types.Int64Value(int64(vrrpInterfaceResp.Priority))
	data.Preempt = types.BoolValue(vrrpInterfaceResp.Preempt)
	data.Master = types.BoolValue(vrrpInterfaceResp.Master)
	data.CreatedAt = types.StringValue(vrrpInterfaceResp.CreatedAt)
	data.UpdatedAt = types.StringValue(vrrpInterfaceResp.UpdatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VRRPInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VRRPInterfaceResourceModel
	var data VRRPInterfaceResourceModel

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

	vrrpInterfaceID := data.ID.ValueString()

	vrrpInterfaceUpdateOpts := vrrpinterfaces.UpdateOpts{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Priority:    int(plan.Priority.ValueInt64()),
		Preempt:     util.ValueKnownBoolPointer(plan.Preempt),
		Master:      util.ValueKnownBoolPointer(plan.Master),
	}

	vrrpInterfaceResp, err := vrrpinterfaces.Update(networkingClient, vrrpInterfaceID, &vrrpinterfaces.VRRPInterfaceUpdate{VRRPInterface: &vrrpInterfaceUpdateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_dc_vrrp_interface", err.Error())
		return
	}

	data.ID = types.StringValue(vrrpInterfaceResp.ID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(vrrpInterfaceResp.Name)
	data.Description = types.StringValue(vrrpInterfaceResp.Description)
	data.DCVRRPID = types.StringValue(vrrpInterfaceResp.DCVRRPID)
	data.DCInterfaceID = types.StringValue(vrrpInterfaceResp.DCInterfaceID)
	data.Priority = types.Int64Value(int64(vrrpInterfaceResp.Priority))
	data.Preempt = types.BoolValue(vrrpInterfaceResp.Preempt)
	data.Master = types.BoolValue(vrrpInterfaceResp.Master)
	data.CreatedAt = types.StringValue(vrrpInterfaceResp.CreatedAt)
	data.UpdatedAt = types.StringValue(vrrpInterfaceResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *VRRPInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VRRPInterfaceResourceModel

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

	err = vrrpinterfaces.Delete(networkingClient, id).ExtractErr()
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_dc_vrrp_interface", err.Error())
		return
	}
}

func (r *VRRPInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
