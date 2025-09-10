package dc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/planmodifiers"
	conntrackhelpers "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dc/v2/conntrack_helpers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

// Ensure the implementation satisfies the desired interfaces.
var (
	_ resource.Resource                = &ConntrackHelperResource{}
	_ resource.ResourceWithConfigure   = &ConntrackHelperResource{}
	_ resource.ResourceWithImportState = &ConntrackHelperResource{}
)

func NewConntrackHelperResource() resource.Resource {
	return &ConntrackHelperResource{}
}

type ConntrackHelperResource struct {
	config clients.Config
}

func (r *ConntrackHelperResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_dc_conntrack_helper"
}

type ConntrackHelperResourceModel struct {
	ID          types.String `tfsdk:"id"`
	DCRouterID  types.String `tfsdk:"dc_router_id"`
	Protocol    types.String `tfsdk:"protocol"`
	Port        types.Int64  `tfsdk:"port"`
	Helper      types.String `tfsdk:"helper"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	Region      types.String `tfsdk:"region"`
}

func (r *ConntrackHelperResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource",
			},

			"dc_router_id": schema.StringAttribute{
				Required:    true,
				Description: "Direct Connect Router ID. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"helper": schema.StringAttribute{
				Required:    true,
				Description: "Helper type. Must be one of: \"ftp\".",
				Validators: []validator.String{
					stringvalidator.OneOf("ftp"),
				},
			},

			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "Protocol. Must be one of: \"tcp\".",
				Validators: []validator.String{
					stringvalidator.OneOf("tcp"),
				},
			},

			"port": schema.Int64Attribute{
				Required:    true,
				Description: "Network port for conntrack target rule.",
			},

			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the conntrack helper",
			},

			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the conntrack helper",
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
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`. Changing this creates a new conntrack_helper.",
			},
		},
		Description: "Manages a direct connect conntrack helper resource._note_ This resource requires Sprut SDN to be enabled in your project.",
	}
}

func (r *ConntrackHelperResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *ConntrackHelperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConntrackHelperResourceModel
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

	conntrackHelperCreateOpts := conntrackhelpers.CreateOpts{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		DCRouterID:  data.DCRouterID.ValueString(),
		Helper:      data.Helper.ValueString(),
		Protocol:    data.Protocol.ValueString(),
		Port:        int(data.Port.ValueInt64()),
	}

	conntrackHelperResp, err := conntrackhelpers.Create(networkingClient, &conntrackhelpers.ConntrackHelperCreate{ConntrackHelper: &conntrackHelperCreateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_dc_conntrack_helper", err.Error())
		return
	}
	conntrackHelperID := conntrackHelperResp.ID
	resp.State.SetAttribute(ctx, path.Root("id"), conntrackHelperID)

	data.ID = types.StringValue(conntrackHelperID)
	data.Region = types.StringValue(region)
	data.DCRouterID = types.StringValue(conntrackHelperResp.DCRouterID)
	data.Helper = types.StringValue(conntrackHelperResp.Helper)
	data.Protocol = types.StringValue(conntrackHelperResp.Protocol)
	data.Port = types.Int64Value(int64(conntrackHelperResp.Port))
	data.Name = types.StringValue(conntrackHelperResp.Name)
	data.Description = types.StringValue(conntrackHelperResp.Description)
	data.CreatedAt = types.StringValue(conntrackHelperResp.CreatedAt)
	data.UpdatedAt = types.StringValue(conntrackHelperResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *ConntrackHelperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConntrackHelperResourceModel

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

	conntrackHelperID := data.ID.ValueString()

	conntrackHelperResp, err := conntrackhelpers.Get(networkingClient, conntrackHelperID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_dc_conntrack_helper", checkDeleted.Error())
		}
		return
	}

	data.ID = types.StringValue(conntrackHelperResp.ID)
	data.Region = types.StringValue(region)
	data.DCRouterID = types.StringValue(conntrackHelperResp.DCRouterID)
	data.Helper = types.StringValue(conntrackHelperResp.Helper)
	data.Protocol = types.StringValue(conntrackHelperResp.Protocol)
	data.Port = types.Int64Value(int64(conntrackHelperResp.Port))
	data.Name = types.StringValue(conntrackHelperResp.Name)
	data.Description = types.StringValue(conntrackHelperResp.Description)
	data.CreatedAt = types.StringValue(conntrackHelperResp.CreatedAt)
	data.UpdatedAt = types.StringValue(conntrackHelperResp.UpdatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConntrackHelperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConntrackHelperResourceModel
	var data ConntrackHelperResourceModel

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

	conntrackHelperID := data.ID.ValueString()

	conntrackHelperUpdateOpts := conntrackhelpers.UpdateOpts{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Helper:      plan.Helper.ValueString(),
		Protocol:    plan.Protocol.ValueString(),
		Port:        int(plan.Port.ValueInt64()),
	}

	conntrackHelperResp, err := conntrackhelpers.Update(networkingClient, conntrackHelperID, &conntrackhelpers.ConntrackHelperUpdate{ConntrackHelper: &conntrackHelperUpdateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_dc_conntrack_helper", err.Error())
		return
	}

	data.ID = types.StringValue(conntrackHelperResp.ID)
	data.Region = types.StringValue(region)
	data.DCRouterID = types.StringValue(conntrackHelperResp.DCRouterID)
	data.Helper = types.StringValue(conntrackHelperResp.Helper)
	data.Protocol = types.StringValue(conntrackHelperResp.Protocol)
	data.Port = types.Int64Value(int64(conntrackHelperResp.Port))
	data.Name = types.StringValue(conntrackHelperResp.Name)
	data.Description = types.StringValue(conntrackHelperResp.Description)
	data.CreatedAt = types.StringValue(conntrackHelperResp.CreatedAt)
	data.UpdatedAt = types.StringValue(conntrackHelperResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *ConntrackHelperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConntrackHelperResourceModel

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

	err = conntrackhelpers.Delete(networkingClient, id).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_dc_conntrack_helper", err.Error())
		return
	}
}

func (r *ConntrackHelperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
