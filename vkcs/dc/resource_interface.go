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
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dc/v2/interfaces"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

// Ensure the implementation satisfies the desired interfaces.
var (
	_ resource.Resource                = &InterfaceResource{}
	_ resource.ResourceWithConfigure   = &InterfaceResource{}
	_ resource.ResourceWithImportState = &InterfaceResource{}
)

func NewInterfaceResource() resource.Resource {
	return &InterfaceResource{}
}

type InterfaceResource struct {
	config clients.Config
}

func (r *InterfaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_dc_interface"
}

type InterfaceResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	DCRouterID         types.String `tfsdk:"dc_router_id"`
	NetworkID          types.String `tfsdk:"network_id"`
	SubnetID           types.String `tfsdk:"subnet_id"`
	BGPAnnounceEnabled types.Bool   `tfsdk:"bgp_announce_enabled"`
	PortID             types.String `tfsdk:"port_id"`
	SDN                types.String `tfsdk:"sdn"`
	IPAddress          types.String `tfsdk:"ip_address"`
	IPNetmask          types.Int64  `tfsdk:"ip_netmask"`
	MACAddress         types.String `tfsdk:"mac_address"`
	MTU                types.Int64  `tfsdk:"mtu"`
	CreatedAt          types.String `tfsdk:"created_at"`
	UpdatedAt          types.String `tfsdk:"updated_at"`

	Region types.String `tfsdk:"region"`
}

func (r *InterfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource",
			},

			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the interface",
			},

			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the interface",
			},

			"dc_router_id": schema.StringAttribute{
				Required:    true,
				Description: "Direct Connect Router ID to attach. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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

			"bgp_announce_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable BGP announce of subnet attached to interface. Default is true",
			},

			"ip_address": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "IP Address of the interface. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},

			"port_id": schema.StringAttribute{
				Computed:    true,
				Description: "Port ID",
			},

			"sdn": schema.StringAttribute{
				Computed:    true,
				Description: "SDN where interface was created",
			},

			"ip_netmask": schema.Int64Attribute{
				Computed:    true,
				Description: "IP Netmask",
			},

			"mac_address": schema.StringAttribute{
				Computed:    true,
				Description: "MAC Address of created interface",
			},

			"mtu": schema.Int64Attribute{
				Computed:    true,
				Description: "MTU",
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
		Description: "Manages a direct connect interface resource._note_ This resource requires Sprut SDN to be enabled in your project.",
	}
}

func (r *InterfaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *InterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InterfaceResourceModel
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

	interfaceCreateOpts := interfaces.CreateOpts{
		Name:               data.Name.ValueString(),
		Description:        data.Description.ValueString(),
		DCRouterID:         data.DCRouterID.ValueString(),
		NetworkID:          data.NetworkID.ValueString(),
		SubnetID:           data.SubnetID.ValueString(),
		IPAddress:          data.IPAddress.ValueString(),
		BGPAnnounceEnabled: util.ValueKnownBoolPointer(data.BGPAnnounceEnabled),
	}

	interfaceResp, err := interfaces.Create(networkingClient, &interfaces.InterfaceCreate{Interface: &interfaceCreateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_dc_interface", err.Error())
		return
	}
	interfaceID := interfaceResp.ID
	resp.State.SetAttribute(ctx, path.Root("id"), interfaceID)

	data.ID = types.StringValue(interfaceResp.ID)
	data.Region = types.StringValue(region)

	data.Name = types.StringValue(interfaceResp.Name)
	data.Description = types.StringValue(interfaceResp.Description)
	data.DCRouterID = types.StringValue(interfaceResp.DCRouterID)
	data.NetworkID = types.StringValue(interfaceResp.NetworkID)
	data.SubnetID = types.StringValue(interfaceResp.SubnetID)
	data.BGPAnnounceEnabled = types.BoolValue(interfaceResp.BGPAnnounceEnabled)
	data.PortID = types.StringValue(interfaceResp.PortID)
	data.SDN = types.StringValue(interfaceResp.SDN)
	data.IPAddress = types.StringValue(interfaceResp.IPAddress)
	data.IPNetmask = types.Int64Value(int64(interfaceResp.IPNetmask))
	data.MACAddress = types.StringValue(interfaceResp.MACAddress)
	data.MTU = types.Int64Value(int64(interfaceResp.MTU))
	data.CreatedAt = types.StringValue(interfaceResp.CreatedAt)
	data.UpdatedAt = types.StringValue(interfaceResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *InterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InterfaceResourceModel

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

	interfaceID := data.ID.ValueString()

	interfaceResp, err := interfaces.Get(networkingClient, interfaceID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_dc_interface", checkDeleted.Error())
		}
		return
	}

	data.ID = types.StringValue(interfaceResp.ID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(interfaceResp.Name)
	data.Description = types.StringValue(interfaceResp.Description)
	data.DCRouterID = types.StringValue(interfaceResp.DCRouterID)
	data.NetworkID = types.StringValue(interfaceResp.NetworkID)
	data.SubnetID = types.StringValue(interfaceResp.SubnetID)
	data.BGPAnnounceEnabled = types.BoolValue(interfaceResp.BGPAnnounceEnabled)
	data.PortID = types.StringValue(interfaceResp.PortID)
	data.SDN = types.StringValue(interfaceResp.SDN)
	data.IPAddress = types.StringValue(interfaceResp.IPAddress)
	data.IPNetmask = types.Int64Value(int64(interfaceResp.IPNetmask))
	data.MACAddress = types.StringValue(interfaceResp.MACAddress)
	data.MTU = types.Int64Value(int64(interfaceResp.MTU))
	data.CreatedAt = types.StringValue(interfaceResp.CreatedAt)
	data.UpdatedAt = types.StringValue(interfaceResp.UpdatedAt)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan InterfaceResourceModel
	var data InterfaceResourceModel

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

	interfaceID := data.ID.ValueString()

	interfaceUpdateOpts := interfaces.UpdateOpts{
		Name:               plan.Name.ValueString(),
		Description:        plan.Description.ValueString(),
		BGPAnnounceEnabled: util.ValueKnownBoolPointer(plan.BGPAnnounceEnabled),
	}

	interfaceResp, err := interfaces.Update(networkingClient, interfaceID, &interfaces.InterfaceUpdate{Interface: &interfaceUpdateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_dc_interface", err.Error())
		return
	}

	data.ID = types.StringValue(interfaceResp.ID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(interfaceResp.Name)
	data.Description = types.StringValue(interfaceResp.Description)
	data.DCRouterID = types.StringValue(interfaceResp.DCRouterID)
	data.NetworkID = types.StringValue(interfaceResp.NetworkID)
	data.SubnetID = types.StringValue(interfaceResp.SubnetID)
	data.BGPAnnounceEnabled = types.BoolValue(interfaceResp.BGPAnnounceEnabled)
	data.PortID = types.StringValue(interfaceResp.PortID)
	data.SDN = types.StringValue(interfaceResp.SDN)
	data.IPAddress = types.StringValue(interfaceResp.IPAddress)
	data.IPNetmask = types.Int64Value(int64(interfaceResp.IPNetmask))
	data.MACAddress = types.StringValue(interfaceResp.MACAddress)
	data.MTU = types.Int64Value(int64(interfaceResp.MTU))
	data.CreatedAt = types.StringValue(interfaceResp.CreatedAt)
	data.UpdatedAt = types.StringValue(interfaceResp.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *InterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InterfaceResourceModel

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

	err = interfaces.Delete(networkingClient, id).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_dc_interface", err.Error())
		return
	}
}

func (r *InterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
