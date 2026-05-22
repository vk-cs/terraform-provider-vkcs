package baremetal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	v1 "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/baremetal/v1"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/baremetal/v1/rents"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/baremetal/v1/servers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

const (
	serverDelay         = 10 * time.Second
	serverMinTimeout    = 10 * time.Second
	serverCreateTimeout = 30 * time.Minute
)

const (
	reprovisionKey = "reprovision"
)

type reprovision struct {
	Enabled bool
}

var (
	_ resource.Resource                     = &ServerResource{}
	_ resource.ResourceWithConfigure        = &ServerResource{}
	_ resource.ResourceWithImportState      = &ServerResource{}
	_ resource.ResourceWithModifyPlan       = &ServerResource{}
	_ resource.ResourceWithConfigValidators = &ServerResource{}
)

func NewServerResource() resource.Resource {
	return &ServerResource{}
}

type ServerResource struct {
	config clients.Config
}

type ServerResourceModel struct {
	ID               types.String   `tfsdk:"id"`
	Name             types.String   `tfsdk:"name"`
	Region           types.String   `tfsdk:"region"`
	AvailabilityZone types.String   `tfsdk:"availability_zone"`
	FlavorID         types.String   `tfsdk:"flavor_id"`
	KeyPair          types.String   `tfsdk:"key_pair"`
	UserData         types.String   `tfsdk:"user_data"`
	OsID             types.String   `tfsdk:"os_id"`
	RaidType         types.String   `tfsdk:"raid_type"`
	Nics             []NicModel     `tfsdk:"nic"`
	Bonds            []BondModel    `tfsdk:"bond"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

type NicModel struct {
	Name  types.String `tfsdk:"name"`
	Vlans []VlanModel  `tfsdk:"vlan"`
}

type BondModel struct {
	Name           types.String `tfsdk:"name"`
	InterfaceNames types.List   `tfsdk:"interface_names"`
	Vlans          []VlanModel  `tfsdk:"vlan"`
}

type VlanModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Native    types.Bool   `tfsdk:"native"`
	NetworkId types.String `tfsdk:"network_id"`
	SubnetId  types.String `tfsdk:"subnet_id"`
}

func (r *ServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_baremetal_server"
}

func (r *ServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the bare metal server.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the bare metal server.",
			},
			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The region to fetch the bare metal server from, defaults to the provider's region.",
			},
			"availability_zone": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Availability zone. If not specified, we will chose the availability zone for you.",
			},
			"flavor_id": schema.StringAttribute{
				Required:    true,
				Description: "Server flavor to rent.",
			},
			"key_pair": schema.StringAttribute{
				Required:    true,
				Description: "The name of a key pair to put on the server. The key pair must already be created and associated with the tenant's account. Changing this creates a new server.",
			},
			"user_data": schema.StringAttribute{
				Optional:    true,
				Description: "Provide the cloud-init user-data payload.",
			},
			"os_id": schema.StringAttribute{
				Optional:    true,
				Description: "Set os id.",
			},
			"raid_type": schema.StringAttribute{
				Optional:    true,
				Description: "Parameter to determine should RAID be used during image flashing.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
		Blocks: map[string]schema.Block{
			"nic": schema.ListNestedBlock{
				Description: "Physical network interfaces.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Interface name (e.g. nic0, eno1). Acts as unique identifier.",
						},
					},
					Blocks: map[string]schema.Block{
						"vlan": schema.ListNestedBlock{
							Description:  "VLAN configuration. Allowed only if interface is not part of a bond.",
							NestedObject: vlanObject(),
						},
					},
				},
			},
			"bond": schema.ListNestedBlock{
				Description: "Link aggregation interfaces (bonds).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Bond interface name (e.g. bond0).",
						},
						"interface_names": schema.ListAttribute{
							Required:    true,
							Description: "List of interface names participating in the bond.",
							ElementType: types.StringType,
						},
					},
					Blocks: map[string]schema.Block{
						"vlan": schema.ListNestedBlock{
							Description:  "VLAN configuration applied to the bond.",
							NestedObject: vlanObject(),
						},
					},
				},
			},
		},
	}
}

func vlanObject() schema.NestedBlockObject {
	return schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of the VLAN.",
			},
			"native": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the VLAN is native.",
			},
			"network_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of the network.",
			},
			"subnet_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of the subnet.",
			},
		},
	}
}

func (r *ServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.config = req.ProviderData.(clients.Config)
}

func (r *ServerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("nic"),
			path.MatchRoot("bond"),
		),
	}
}

func (r *ServerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	reprovisionState, err := json.Marshal(reprovision{Enabled: true})
	if err != nil {
		resp.Diagnostics.AddError("Error encoding reprovision state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.Private.SetKey(ctx, reprovisionKey, reprovisionState)...)

	var plan ServerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := plan.Region.ValueString()
	if plan.Region.IsUnknown() {
		region = r.config.GetRegion()
		plan.Region = types.StringValue(region)
		resp.Plan.SetAttribute(ctx, path.Root("region"), region)
	}

	if req.State.Raw.IsNull() {
		return
	}

	var state ServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Region.ValueString() != region {
		return
	}

	if !plan.FlavorID.Equal(state.FlavorID) {
		resp.Diagnostics.AddError(
			"Field cannot be updated",
			"flavor_id cannot be changed after creation",
		)
	}

	if !plan.AvailabilityZone.Equal(state.AvailabilityZone) {
		resp.Diagnostics.AddError(
			"Field cannot be updated",
			"availability_zone cannot be changed after creation",
		)
	}
}

func (r *ServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	ctx = tflog.SetField(ctx, "region", region)
	client, err := r.config.BareMetalV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS baremetal API client", err.Error())
		return
	}

	serverID, d := rent(ctx, data, client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = types.StringValue(serverID)

	renameOpts := servers.RenameOpts{
		ServerName: data.Name.ValueString(),
	}

	if err := servers.Rename(client, data.ID.ValueString(), &renameOpts).ExtractErr(); err != nil {
		resp.Diagnostics.AddError("Error renaming server", err.Error())
		return
	}

	data.Region = types.StringValue(region)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := state.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.BareMetalV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS baremetal API client", err.Error())
		return
	}

	serverID := state.ID.ValueString()
	if serverID == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	ctx = tflog.SetField(ctx, "server_id", serverID)

	tflog.Debug(ctx, "Calling VKCS baremetal API to retrieve server by id", map[string]interface{}{"id": serverID})
	server, err := servers.Get(client, serverID).Extract()
	if errutil.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading baremetal server", err.Error())
		return
	}

	state.Name = types.StringValue(server.ServerName)
	state.Region = types.StringValue(region)
	state.AvailabilityZone = types.StringValue(server.AvailabilityZone)
	state.FlavorID = types.StringPointerValue(server.FlavorId)
	state.RaidType = types.StringPointerValue(server.RaidType)
	state.OsID = types.StringPointerValue(server.ImageId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServerResourceModel
	var state ServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := state.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.BareMetalV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS baremetal API client", err.Error())
		return
	}

	serverID := state.ID.ValueString()
	resp.Diagnostics.Append(provision(ctx, plan, client)...)

	renameOpts := servers.RenameOpts{
		ServerName: plan.Name.ValueString(),
	}

	tflog.Debug(ctx, "Calling VKCS baremetal API to update server by id", map[string]interface{}{"id": serverID})
	if err := servers.Rename(client, serverID, &renameOpts).ExtractErr(); err != nil {
		resp.Diagnostics.AddError("Error updating baremetal server", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	val, diags := req.Private.GetKey(ctx, reprovisionKey)
	resp.Diagnostics.Append(diags...)
	if val != nil {
		var reprovisionState reprovision
		if err := json.Unmarshal(val, &reprovisionState); err != nil {
			resp.Diagnostics.AddError("Error encoding reprovision state", err.Error())
			return
		}

		if reprovisionState.Enabled {
			tflog.Info(ctx, "Server will be reprovisioned")
			return
		}
	}

	var state ServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := state.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.BareMetalV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS baremetal API client", err.Error())
		return
	}

	serverID := state.ID.ValueString()
	if serverID == "" {
		return
	}

	ctx = tflog.SetField(ctx, "server_id", serverID)
	tflog.Debug(ctx, "Calling VKCS baremetal API to delete server by id", map[string]interface{}{"id": serverID})
	if err := servers.Delete(client, serverID).ExtractErr(); err != nil && !errutil.IsNotFound(err) {
		resp.Diagnostics.AddError("Error deleting baremetal server", err.Error())
		return
	}
}

func (r *ServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func rent(ctx context.Context, data ServerResourceModel, client *gophercloud.ServiceClient) (serverID string, diags diag.Diagnostics) {
	timeout, d := data.Timeouts.Create(ctx, serverCreateTimeout)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	createOpts := rents.CreateOpts{
		ProvisionFields:  defineProvisionFields(ctx, data),
		FlavorId:         data.FlavorID.ValueString(),
		ServerCount:      1,
		AvailabilityZone: data.AvailabilityZone.ValueStringPointer(),
	}

	tflog.Debug(ctx, "Calling VKCS baremetal API to create rent request")
	rentReq, err := rents.Create(client, &createOpts).Extract()
	if err != nil {
		diags.AddError("Error creating baremetal server", err.Error())
		return
	}

	if len(rentReq.ServerIds) != 1 {
		diags.AddError("Error retrieving baremetal servers ID from list", fmt.Sprintf("%+v", rentReq.ServerIds))
		return
	}

	serverID = rentReq.ServerIds[0]
	tflog.SetField(ctx, "rent_request_id", rentReq.RentRequestId)

	serverStateConf := &retry.StateChangeConf{
		Pending:    []string{string(servers.StatusDiscovered), string(servers.StatusInProgress)},
		Target:     []string{string(servers.StatusActive)},
		Refresh:    provisionRefreshFunc(ctx, client, serverID),
		Timeout:    timeout,
		Delay:      serverDelay,
		MinTimeout: serverMinTimeout,
	}

	if _, err := serverStateConf.WaitForStateContext(ctx); err != nil {
		diags.AddError("Error waiting baremetal server", err.Error())
		return
	}

	tflog.Debug(ctx, "Renting baremetal done")
	return
}

func defineProvisionFields(ctx context.Context, data ServerResourceModel) v1.ProvisionFields {
	fields := v1.ProvisionFields{
		ProvisionType:     v1.ProvisionTypeNOOS,
		RaidType:          data.RaidType.ValueStringPointer(),
		KeypairName:       data.KeyPair.ValueString(),
		UserData:          encodeUserData(data.UserData), // base64 encode
		NetworkInterfaces: flattenNetworkInterfaces(data.Nics),
		Bonds:             flattenBonds(ctx, data.Bonds),
	}

	if !data.OsID.IsNull() && !data.OsID.IsUnknown() {
		fields.ProvisionType = v1.ProvisionTypeIMAGE
		fields.ImageSource = v1.ImageSourcePUBLIC
		fields.ImageId = data.OsID.ValueStringPointer()
	}

	return fields
}

func provision(ctx context.Context, data ServerResourceModel, client *gophercloud.ServiceClient) (diags diag.Diagnostics) {
	timeout, d := data.Timeouts.Create(ctx, serverCreateTimeout)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	provisionOpts := servers.ProvisionOpts{
		ProvisionFields: defineProvisionFields(ctx, data),
	}

	serverID := data.ID.ValueString()

	tflog.Debug(ctx, "Calling VKCS baremetal API to provisioning server by id", map[string]interface{}{"id": serverID})
	if err := servers.Provision(client, serverID, &provisionOpts).ExtractErr(); err != nil {
		diags.AddError("Error provisioning baremetal server", err.Error())
		return
	}

	serverStateConf := &retry.StateChangeConf{
		Pending:    []string{string(servers.StatusDiscovered), string(servers.StatusInProgress)},
		Target:     []string{string(servers.StatusActive)},
		Refresh:    provisionRefreshFunc(ctx, client, serverID),
		Timeout:    timeout,
		Delay:      serverDelay,
		MinTimeout: serverMinTimeout,
	}

	if _, err := serverStateConf.WaitForStateContext(ctx); err != nil {
		diags.AddError("Error waiting baremetal server", err.Error())
		return
	}

	tflog.Debug(ctx, "Renting baremetal done")
	return
}

func provisionRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		tflog.Debug(ctx, "Calling VKCS baremetal API to retrieve server by id", map[string]interface{}{"id": id})
		server, err := servers.Get(client, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return server, string(server.Status), nil
	}
}

func flattenNetworkInterfaces(items []NicModel) []*v1.NetworkInterfaceConfig {
	configs := make([]*v1.NetworkInterfaceConfig, 0, len(items))

	for _, item := range items {
		configs = append(configs, &v1.NetworkInterfaceConfig{
			NicName: item.Name.ValueString(),
			Vlans:   flattenVlans(item.Vlans),
		})
	}

	sort.SliceStable(configs, func(i, j int) bool {
		return configs[i].NicName < configs[j].NicName
	})

	return configs
}

func flattenVlans(items []VlanModel) []*v1.VlanConfig {
	vlans := make([]*v1.VlanConfig, 0, len(items))

	for _, vlan := range items {
		vlans = append(vlans, &v1.VlanConfig{
			VlanId:    vlan.ID.ValueInt64Pointer(),
			IsNative:  vlan.Native.ValueBool(),
			NetworkId: vlan.NetworkId.ValueString(),
			SubnetId:  vlan.SubnetId.ValueString(),
		})
	}

	sort.SliceStable(vlans, func(i, j int) bool {
		vi, vj := int64(0), int64(0)

		if vlans[i].VlanId != nil {
			vi = *vlans[i].VlanId
		}
		if vlans[j].VlanId != nil {
			vj = *vlans[j].VlanId
		}
		return vi < vj
	})

	return vlans
}

func flattenBonds(ctx context.Context, items []BondModel) []*v1.BondConfig {
	bonds := make([]*v1.BondConfig, 0, len(items))
	for _, item := range items {
		var ifNames []string
		item.InterfaceNames.ElementsAs(ctx, &ifNames, false)

		bonds = append(bonds, &v1.BondConfig{
			BondName:       item.Name.ValueString(),
			InterfaceNames: ifNames,
			Vlans:          flattenVlans(item.Vlans),
		})
	}

	sort.SliceStable(bonds, func(i, j int) bool {
		return bonds[i].BondName < bonds[j].BondName
	})

	return bonds
}

func encodeUserData(data types.String) *string {
	if data.IsNull() {
		return nil
	}
	value := base64.URLEncoding.EncodeToString([]byte(data.ValueString()))

	return &value
}
