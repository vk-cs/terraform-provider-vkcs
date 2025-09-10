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
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dc/v2/staticroutes"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

// Ensure the implementation satisfies the desired interfaces.
var (
	_ resource.Resource                = &StaticRouteResource{}
	_ resource.ResourceWithConfigure   = &StaticRouteResource{}
	_ resource.ResourceWithImportState = &StaticRouteResource{}
)

func NewStaticRouteResource() resource.Resource {
	return &StaticRouteResource{}
}

type StaticRouteResource struct {
	config clients.Config
}

func (r *StaticRouteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_dc_static_route"
}

type StaticRouteResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	DCRouterID  types.String `tfsdk:"dc_router_id"`
	Network     types.String `tfsdk:"network"`
	Gateway     types.String `tfsdk:"gateway"`
	Metric      types.Int64  `tfsdk:"metric"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	Region      types.String `tfsdk:"region"`
}

func (r *StaticRouteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource",
			},

			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the static route",
			},

			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the static route",
			},

			"dc_router_id": schema.StringAttribute{
				Required:    true,
				Description: "Direct Connect Router ID to attach. Changing this creates a new resource",
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

			"metric": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Metric to use for route. Default is 1",
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
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`. Changing this creates a new static_route.",
			},
		},
		Description: "Manages a direct connect BGP Static Announce resource._note_ This resource requires Sprut SDN to be enabled in your project.",
	}
}

func (r *StaticRouteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *StaticRouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data StaticRouteResourceModel
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

	staticRouteCreateOpts := staticroutes.CreateOpts{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		DCRouterID:  data.DCRouterID.ValueString(),
		Network:     data.Network.ValueString(),
		Gateway:     data.Gateway.ValueString(),
		Metric:      int(data.Metric.ValueInt64()),
	}

	staticRouteResp, err := staticroutes.Create(networkingClient, &staticroutes.StaticRouteCreate{StaticRoute: &staticRouteCreateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_dc_static_route", err.Error())
		return
	}
	staticRouteID := staticRouteResp.ID
	resp.State.SetAttribute(ctx, path.Root("id"), staticRouteID)

	data.ID = types.StringValue(staticRouteID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(staticRouteResp.Name)
	data.Description = types.StringValue(staticRouteResp.Description)
	data.DCRouterID = types.StringValue(staticRouteResp.DCRouterID)
	data.Network = types.StringValue(staticRouteResp.Network)
	data.Gateway = types.StringValue(staticRouteResp.Gateway)
	data.Metric = types.Int64Value(int64(staticRouteResp.Metric))
	data.CreatedAt = types.StringValue(staticRouteResp.CreatedAt)
	data.UpdatedAt = types.StringValue(staticRouteResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *StaticRouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data StaticRouteResourceModel

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

	staticRouteID := data.ID.ValueString()

	staticRouteResp, err := staticroutes.Get(networkingClient, staticRouteID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_dc_static_route", checkDeleted.Error())
		}
		return
	}

	data.ID = types.StringValue(staticRouteID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(staticRouteResp.Name)
	data.Description = types.StringValue(staticRouteResp.Description)
	data.DCRouterID = types.StringValue(staticRouteResp.DCRouterID)
	data.Network = types.StringValue(staticRouteResp.Network)
	data.Gateway = types.StringValue(staticRouteResp.Gateway)
	data.Metric = types.Int64Value(int64(staticRouteResp.Metric))
	data.CreatedAt = types.StringValue(staticRouteResp.CreatedAt)
	data.UpdatedAt = types.StringValue(staticRouteResp.UpdatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StaticRouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StaticRouteResourceModel
	var data StaticRouteResourceModel

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

	staticRouteID := data.ID.ValueString()

	staticRouteUpdateOpts := staticroutes.UpdateOpts{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Metric:      int(plan.Metric.ValueInt64()),
	}

	staticRouteResp, err := staticroutes.Update(networkingClient, staticRouteID, &staticroutes.StaticRouteUpdate{StaticRoute: &staticRouteUpdateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_dc_static_route", err.Error())
		return
	}

	data.ID = types.StringValue(staticRouteID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(staticRouteResp.Name)
	data.Description = types.StringValue(staticRouteResp.Description)
	data.DCRouterID = types.StringValue(staticRouteResp.DCRouterID)
	data.Network = types.StringValue(staticRouteResp.Network)
	data.Gateway = types.StringValue(staticRouteResp.Gateway)
	data.Metric = types.Int64Value(int64(staticRouteResp.Metric))
	data.CreatedAt = types.StringValue(staticRouteResp.CreatedAt)
	data.UpdatedAt = types.StringValue(staticRouteResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *StaticRouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data StaticRouteResourceModel

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

	err = staticroutes.Delete(networkingClient, id).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_dc_static_route", err.Error())
		return
	}
}

func (r *StaticRouteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
