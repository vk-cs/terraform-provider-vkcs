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
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dc/v2/ipportforwardings"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

// Ensure the implementation satisfies the desired interfaces.
var (
	_ resource.Resource                = &IPPortForwardingResource{}
	_ resource.ResourceWithConfigure   = &IPPortForwardingResource{}
	_ resource.ResourceWithImportState = &IPPortForwardingResource{}
)

func NewIPPortForwardingResource() resource.Resource {
	return &IPPortForwardingResource{}
}

type IPPortForwardingResource struct {
	config clients.Config
}

func (r *IPPortForwardingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_dc_ip_port_forwarding"
}

type IPPortForwardingResourceModel struct {
	ID            types.String `tfsdk:"id"`
	DCInterfaceID types.String `tfsdk:"dc_interface_id"`
	Protocol      types.String `tfsdk:"protocol"`
	Source        types.String `tfsdk:"source"`
	Destination   types.String `tfsdk:"destination"`
	Port          types.Int64  `tfsdk:"port"`
	ToDestination types.String `tfsdk:"to_destination"`
	ToPort        types.Int64  `tfsdk:"to_port"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	Region        types.String `tfsdk:"region"`
}

func (r *IPPortForwardingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource",
			},

			"dc_interface_id": schema.StringAttribute{
				Required:    true,
				Description: "Direct Connect Interface ID. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"to_destination": schema.StringAttribute{
				Required:    true,
				Description: "IP Address of forwarding's destination.",
			},

			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "Protocol. Must be one of: \"tcp\", \"udp\", \"any\".",
				Validators: []validator.String{
					stringvalidator.OneOf("tcp", "udp", "any"),
				},
			},

			"source": schema.StringAttribute{
				Optional:    true,
				Description: "Source address selector.",
			},

			"destination": schema.StringAttribute{
				Optional:    true,
				Description: "Destination address selector.",
			},

			"port": schema.Int64Attribute{
				Optional:    true,
				Description: "Port selector.",
			},

			"to_port": schema.Int64Attribute{
				Optional:    true,
				Description: "Destination port selector.",
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
				Optional:    true,
				Computed:    true,
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`.",
			},
		},
		Description: "Manages a direct connect ip port forwarding resource._note_ This resource requires Sprut SDN to be enabled in your project.",
	}
}

func (r *IPPortForwardingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *IPPortForwardingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IPPortForwardingResourceModel
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

	ipPortForwardingCreateOpts := ipportforwardings.CreateOpts{
		Name:          data.Name.ValueString(),
		Description:   data.Description.ValueString(),
		DCInterfaceID: data.DCInterfaceID.ValueString(),
		Protocol:      data.Protocol.ValueString(),
		ToDestination: data.ToDestination.ValueString(),
		Source:        data.Source.ValueStringPointer(),
		Destination:   data.Destination.ValueStringPointer(),
		Port:          data.Port.ValueInt64Pointer(),
		ToPort:        data.ToPort.ValueInt64Pointer(),
	}

	ipPortForwardingResp, err := ipportforwardings.Create(networkingClient, &ipportforwardings.IPPortForwardingCreate{IPPortForwarding: &ipPortForwardingCreateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_dc_ip_port_forwarding", err.Error())
		return
	}
	ipPortForwardingID := ipPortForwardingResp.ID
	resp.State.SetAttribute(ctx, path.Root("id"), ipPortForwardingID)

	data.ID = types.StringValue(ipPortForwardingID)
	data.Region = types.StringValue(region)
	data.DCInterfaceID = types.StringValue(ipPortForwardingResp.DCInterfaceID)
	data.Protocol = types.StringValue(ipPortForwardingResp.Protocol)
	data.ToDestination = types.StringValue(ipPortForwardingResp.ToDestination)
	data.Name = types.StringValue(ipPortForwardingResp.Name)
	data.Description = types.StringValue(ipPortForwardingResp.Description)
	data.CreatedAt = types.StringValue(ipPortForwardingResp.CreatedAt)
	data.UpdatedAt = types.StringValue(ipPortForwardingResp.UpdatedAt)
	data.Source = types.StringPointerValue(ipPortForwardingResp.Source)
	data.Destination = types.StringPointerValue(ipPortForwardingResp.Destination)
	data.Port = types.Int64PointerValue(ipPortForwardingResp.Port)
	data.ToPort = types.Int64PointerValue(ipPortForwardingResp.ToPort)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *IPPortForwardingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IPPortForwardingResourceModel

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

	ipPortForwardingID := data.ID.ValueString()

	ipPortForwardingResp, err := ipportforwardings.Get(networkingClient, ipPortForwardingID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_dc_ip_port_forwarding", checkDeleted.Error())
		}
		return
	}

	data.ID = types.StringValue(ipPortForwardingID)
	data.Region = types.StringValue(region)
	data.DCInterfaceID = types.StringValue(ipPortForwardingResp.DCInterfaceID)
	data.Protocol = types.StringValue(ipPortForwardingResp.Protocol)
	data.ToDestination = types.StringValue(ipPortForwardingResp.ToDestination)
	data.Name = types.StringValue(ipPortForwardingResp.Name)
	data.Description = types.StringValue(ipPortForwardingResp.Description)
	data.CreatedAt = types.StringValue(ipPortForwardingResp.CreatedAt)
	data.UpdatedAt = types.StringValue(ipPortForwardingResp.UpdatedAt)
	data.Source = types.StringPointerValue(ipPortForwardingResp.Source)
	data.Destination = types.StringPointerValue(ipPortForwardingResp.Destination)
	data.Port = types.Int64PointerValue(ipPortForwardingResp.Port)
	data.ToPort = types.Int64PointerValue(ipPortForwardingResp.ToPort)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IPPortForwardingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IPPortForwardingResourceModel
	var data IPPortForwardingResourceModel

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

	ipPortForwardingID := data.ID.ValueString()

	ipPortForwardingUpdateOpts := ipportforwardings.UpdateOpts{
		Name:          plan.Name.ValueString(),
		Description:   plan.Description.ValueString(),
		Protocol:      plan.Protocol.ValueString(),
		ToDestination: plan.ToDestination.ValueString(),
		Source:        plan.Source.ValueStringPointer(),
		Destination:   plan.Destination.ValueStringPointer(),
		Port:          plan.Port.ValueInt64Pointer(),
		ToPort:        plan.ToPort.ValueInt64Pointer(),
	}

	ipPortForwardingResp, err := ipportforwardings.Update(networkingClient, ipPortForwardingID, &ipportforwardings.IPPortForwardingUpdate{IPPortForwarding: &ipPortForwardingUpdateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_dc_ip_port_forwarding", err.Error())
		return
	}

	data.ID = types.StringValue(ipPortForwardingID)
	data.Region = types.StringValue(region)
	data.DCInterfaceID = types.StringValue(ipPortForwardingResp.DCInterfaceID)
	data.Protocol = types.StringValue(ipPortForwardingResp.Protocol)
	data.ToDestination = types.StringValue(ipPortForwardingResp.ToDestination)
	data.Name = types.StringValue(ipPortForwardingResp.Name)
	data.Description = types.StringValue(ipPortForwardingResp.Description)
	data.CreatedAt = types.StringValue(ipPortForwardingResp.CreatedAt)
	data.UpdatedAt = types.StringValue(ipPortForwardingResp.UpdatedAt)
	data.Source = types.StringPointerValue(ipPortForwardingResp.Source)
	data.Destination = types.StringPointerValue(ipPortForwardingResp.Destination)
	data.Port = types.Int64PointerValue(ipPortForwardingResp.Port)
	data.ToPort = types.Int64PointerValue(ipPortForwardingResp.ToPort)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *IPPortForwardingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IPPortForwardingResourceModel

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

	err = ipportforwardings.Delete(networkingClient, id).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_dc_conntrack_helper", err.Error())
		return
	}
}

func (r *IPPortForwardingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
