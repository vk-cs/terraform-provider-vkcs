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
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dc/v2/routers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

// Ensure the implementation satisfies the desired interfaces.
var _ resource.Resource = &RouterResource{}

func NewRouterResource() resource.Resource {
	return &RouterResource{}
}

type RouterResource struct {
	config clients.Config
}

func (r *RouterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_dc_router"
}

type RouterResourceModel struct {
	ID               types.String `tfsdk:"id"`
	AvailabilityZone types.String `tfsdk:"availability_zone"`
	Flavor           types.String `tfsdk:"flavor"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
	Region           types.String `tfsdk:"region"`
}

func (r *RouterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource",
			},

			"availability_zone": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The availability zone in which to create the router. Changing this creates a new router",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"flavor": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Flavor of the router. Possible values can be obtained with vkcs_dc_api_options data source. Changing this creates a new router. _note_ Not to be confused with compute service flavors.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
		Description: "Manages a direct connect router resource. <br> ~> **Note:** This resource requires Sprut SDN to be enabled in your project.",
	}
}

func (r *RouterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *RouterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RouterResourceModel
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

	routerCreateOpts := routers.CreateOpts{
		AvailabilityZone: data.AvailabilityZone.ValueString(),
		Flavor:           data.Flavor.ValueString(),
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueString(),
	}

	routerResp, err := routers.Create(networkingClient, &routers.RouterCreate{Router: &routerCreateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_dc_router", err.Error())
		return
	}
	routerID := routerResp.ID
	resp.State.SetAttribute(ctx, path.Root("id"), routerID)

	data.ID = types.StringValue(routerID)
	data.Region = types.StringValue(region)
	data.AvailabilityZone = types.StringValue(routerResp.AvailabilityZone)
	data.Flavor = types.StringValue(routerResp.Flavor)
	data.Name = types.StringValue(routerResp.Name)
	data.Description = types.StringValue(routerResp.Description)
	data.CreatedAt = types.StringValue(routerResp.CreatedAt)
	data.UpdatedAt = types.StringValue(routerResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *RouterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RouterResourceModel

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

	routerID := data.ID.ValueString()

	routerResp, err := routers.Get(networkingClient, routerID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_dc_router", checkDeleted.Error())
		}
		return
	}
	data.Region = types.StringValue(region)
	data.AvailabilityZone = types.StringValue(routerResp.AvailabilityZone)
	data.Flavor = types.StringValue(routerResp.Flavor)
	data.Name = types.StringValue(routerResp.Name)
	data.Description = types.StringValue(routerResp.Description)
	data.CreatedAt = types.StringValue(routerResp.CreatedAt)
	data.UpdatedAt = types.StringValue(routerResp.UpdatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RouterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RouterResourceModel
	var data RouterResourceModel

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

	routerID := data.ID.ValueString()

	routerUpdateOpts := routers.UpdateOpts{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	routerResp, err := routers.Update(networkingClient, routerID, &routers.RouterUpdate{Router: &routerUpdateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_dc_router", err.Error())
		return
	}

	data.ID = types.StringValue(routerResp.ID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(routerResp.Name)
	data.Description = types.StringValue(routerResp.Description)
	data.CreatedAt = types.StringValue(routerResp.CreatedAt)
	data.UpdatedAt = types.StringValue(routerResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *RouterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RouterResourceModel

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

	err = routers.Delete(networkingClient, id).ExtractErr()
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_dc_router", err.Error())
		return
	}
}

func (r *RouterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
