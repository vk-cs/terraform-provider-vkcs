package baremetal

import (
	"cmp"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strings"
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
	ID               types.String        `tfsdk:"id"`
	Name             types.String        `tfsdk:"name"`
	Region           types.String        `tfsdk:"region"`
	AvailabilityZone types.String        `tfsdk:"availability_zone"`
	FlavorID         types.String        `tfsdk:"flavor_id"`
	KeyPair          types.String        `tfsdk:"key_pair"`
	UserData         types.String        `tfsdk:"user_data"`
	OsID             types.String        `tfsdk:"os_id"`
	Monitoring       types.Bool          `tfsdk:"monitoring"`
	StorageLayout    *StorageLayoutModel `tfsdk:"storage_layout"`
	Nics             []NicModel          `tfsdk:"nic"`
	Bonds            []BondModel         `tfsdk:"bond"`
	Timeouts         timeouts.Value      `tfsdk:"timeouts"`
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

type StorageLayoutModel struct {
	Disks []StorageDiskModel `tfsdk:"disk"`
	Raids []StorageRaidModel `tfsdk:"raid"`
}

type StorageDiskModel struct {
	Id         types.String            `tfsdk:"id"`
	Type       types.String            `tfsdk:"type"`
	Size       types.Int64             `tfsdk:"size"`
	Partitions []StoragePartitionModel `tfsdk:"partition"`
}

type StorageRaidModel struct {
	Id         types.String            `tfsdk:"id"`
	Type       types.String            `tfsdk:"type"`
	Members    types.List              `tfsdk:"members"`
	Partitions []StoragePartitionModel `tfsdk:"partition"`
}

// StoragePartitionModel carries its own mount and filesystem type. A null or
// empty mount means the partition is created but not mounted; both
// representations round-trip through the API as supplied.
type StoragePartitionModel struct {
	Mount types.String `tfsdk:"mount"`
	Fs    types.String `tfsdk:"fs"`
	Size  types.String `tfsdk:"size"`
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
			"monitoring": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the monitoring is actively enabled.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
		Blocks: map[string]schema.Block{
			"storage_layout": schema.SingleNestedBlock{
				Description: "Storage layout of the bare metal server: disks carry their own partitions, " +
					"raids are assembled from whole disks. Changing this triggers reprovisioning.",
				Blocks: map[string]schema.Block{
					"disk": schema.ListNestedBlock{
						Description: "Logical disks and their partition layout.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Required:    true,
									Description: "Logical disk identifier.",
								},
								"type": schema.StringAttribute{
									Required:    true,
									Description: "Storage medium of the disk: SSD, HDD or NVME (case-insensitive). Must match the flavor disk type.",
								},
								"size": schema.Int64Attribute{
									Required:    true,
									Description: "Declared disk size in whole GiB, taken from the flavor. Used to pick a real disk within the size tolerance.",
								},
							},
							Blocks: map[string]schema.Block{
								"partition": partitionBlock(),
							},
						},
					},
					"raid": schema.ListNestedBlock{
						Description: "RAID arrays assembled from whole disks; partitions are cut on top of the md device.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Required:    true,
									Description: "RAID identifier.",
								},
								"type": schema.StringAttribute{
									Required:    true,
									Description: "RAID type: raid1 (case-insensitive).",
								},
								"members": schema.ListAttribute{
									Required:    true,
									ElementType: types.StringType,
									Description: "Disk identifiers the RAID is assembled from. All members must share the type and the declared size.",
								},
							},
							Blocks: map[string]schema.Block{
								"partition": partitionBlock(),
							},
						},
					},
				},
			},
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

	// Enum-like layout values are case-insensitive: normalize the plan to the
	// canonical lowercase form so any spelling compares equal against the
	// state (flatten canonicalizes the same way).
	if plan.StorageLayout != nil {
		plan.StorageLayout = flattenStorageLayout(ctx, expandStorageLayout(ctx, plan.StorageLayout))
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("storage_layout"), plan.StorageLayout)...)
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

	updateOpts := servers.UpdateOpts{
		ServerName: data.Name.ValueString(),
	}

	if err := servers.Update(client, data.ID.ValueString(), &updateOpts).ExtractErr(); err != nil {
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
	state.Monitoring = types.BoolPointerValue(server.Monitoring)
	state.OsID = types.StringPointerValue(server.ImageId)
	state.StorageLayout = flattenStorageLayout(ctx, server.StorageLayout)

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
	plan.ID = state.ID

	act := action(plan, state)
	if act != nil {
		resp.Diagnostics.Append(act(ctx, plan, client)...)
	}

	updateOpts := servers.UpdateOpts{
		ServerName: plan.Name.ValueString(),
	}

	tflog.Debug(ctx, "Calling VKCS baremetal API to update server by id", map[string]interface{}{"id": serverID})
	if err := servers.Update(client, serverID, &updateOpts).ExtractErr(); err != nil {
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

func action(plan, state ServerResourceModel) func(context.Context, ServerResourceModel, *gophercloud.ServiceClient) diag.Diagnostics {
	if needProvision(plan, state) {
		return provision
	}
	if needUpdateNetworkConfig(plan, state) {
		return updateNetworkConfig
	}

	return nil
}

func needProvision(plan, state ServerResourceModel) bool {
	return !plan.OsID.Equal(state.OsID) ||
		!plan.UserData.Equal(state.UserData) ||
		!plan.KeyPair.Equal(state.KeyPair) ||
		!plan.Monitoring.Equal(state.Monitoring) ||
		!equalStorageLayout(plan.StorageLayout, state.StorageLayout)
}

func needUpdateNetworkConfig(plan, state ServerResourceModel) bool {
	return !equalNICs(plan.Nics, state.Nics) || !equalBonds(plan.Bonds, state.Bonds)
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

	tflog.Debug(ctx, "Provisioning baremetal done")
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

func updateNetworkConfig(ctx context.Context, data ServerResourceModel, client *gophercloud.ServiceClient) (diags diag.Diagnostics) {
	timeout, d := data.Timeouts.Create(ctx, serverCreateTimeout)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	updateNetworkConfigOpts := servers.UpdateNetworkConfigOpts{
		NetworkInterfaces: flattenNetworkInterfaces(data.Nics),
		Bonds:             flattenBonds(ctx, data.Bonds),
	}

	serverID := data.ID.ValueString()

	tflog.Debug(ctx, "Calling VKCS baremetal API to updating network config by server id", map[string]interface{}{"id": serverID})
	if err := servers.UpdateNetworkConfig(client, serverID, &updateNetworkConfigOpts).ExtractErr(); err != nil {
		diags.AddError("Error update network config for baremetal server", err.Error())
		return
	}

	serverStateConf := &retry.StateChangeConf{
		Pending:    []string{string(servers.StatusInProgress)},
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

	tflog.Debug(ctx, "Updating network config done")
	return
}

func defineProvisionFields(ctx context.Context, data ServerResourceModel) v1.ProvisionFields {
	fields := v1.ProvisionFields{
		ProvisionType:     v1.ProvisionTypeNOOS,
		KeypairName:       data.KeyPair.ValueString(),
		UserData:          encodeUserData(data.UserData), // base64 encode
		Monitoring:        data.Monitoring.ValueBoolPointer(),
		NetworkInterfaces: flattenNetworkInterfaces(data.Nics),
		Bonds:             flattenBonds(ctx, data.Bonds),
		StorageLayout:     expandStorageLayout(ctx, data.StorageLayout),
	}

	if !data.OsID.IsNull() && !data.OsID.IsUnknown() {
		fields.ProvisionType = v1.ProvisionTypeIMAGE
		fields.ImageSource = v1.ImageSourcePUBLIC
		fields.ImageId = data.OsID.ValueStringPointer()
	}

	return fields
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

func equalNICs(a, b []NicModel) bool {
	if len(a) != len(b) {
		return false
	}

	a = slices.Clone(a)
	b = slices.Clone(b)

	slices.SortFunc(a, func(a, b NicModel) int {
		return cmp.Compare(
			a.Name.ValueString(),
			b.Name.ValueString(),
		)
	})

	slices.SortFunc(b, func(a, b NicModel) int {
		return cmp.Compare(
			a.Name.ValueString(),
			b.Name.ValueString(),
		)
	})

	return slices.EqualFunc(a, b, func(a, b NicModel) bool {
		return a.Name.ValueString() == b.Name.ValueString() &&
			equalVlans(a.Vlans, b.Vlans)
	})
}

func equalBonds(a, b []BondModel) bool {
	if len(a) != len(b) {
		return false
	}

	a = slices.Clone(a)
	b = slices.Clone(b)

	slices.SortFunc(a, func(a, b BondModel) int {
		return cmp.Compare(
			a.Name.ValueString(),
			b.Name.ValueString(),
		)
	})

	slices.SortFunc(b, func(a, b BondModel) int {
		return cmp.Compare(
			a.Name.ValueString(),
			b.Name.ValueString(),
		)
	})

	return slices.EqualFunc(a, b, func(a, b BondModel) bool {
		return a.Name.ValueString() == b.Name.ValueString() &&
			equalStringLists(a.InterfaceNames, b.InterfaceNames) &&
			equalVlans(a.Vlans, b.Vlans)
	})
}

func equalVlans(a, b []VlanModel) bool {
	if len(a) != len(b) {
		return false
	}

	a = slices.Clone(a)
	b = slices.Clone(b)

	slices.SortFunc(a, func(a, b VlanModel) int {
		return cmp.Compare(
			a.ID.ValueInt64(),
			b.ID.ValueInt64(),
		)
	})

	slices.SortFunc(b, func(a, b VlanModel) int {
		return cmp.Compare(
			a.ID.ValueInt64(),
			b.ID.ValueInt64(),
		)
	})

	return slices.EqualFunc(a, b, func(a, b VlanModel) bool {
		return a.ID.ValueInt64() == b.ID.ValueInt64() &&
			a.Native.ValueBool() == b.Native.ValueBool() &&
			a.NetworkId.ValueString() == b.NetworkId.ValueString() &&
			a.SubnetId.ValueString() == b.SubnetId.ValueString()
	})
}

func equalStringLists(a, b types.List) bool {
	var aValues []string
	var bValues []string

	a.ElementsAs(context.Background(), &aValues, false)
	b.ElementsAs(context.Background(), &bValues, false)

	slices.Sort(aValues)
	slices.Sort(bValues)

	return slices.Equal(aValues, bValues)
}

func partitionBlock() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Ordered partitions of the device; order determines placement on disk.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"mount": schema.StringAttribute{
					Optional: true,
					Description: "Mount point of the partition. Empty or omitted means the partition is " +
						"created but not mounted.",
				},
				"fs": schema.StringAttribute{
					Required:    true,
					Description: "Filesystem type: ext4, xfs, vfat or swap (case-insensitive). Swap requires an empty mount.",
				},
				"size": schema.StringAttribute{
					Required:    true,
					Description: "Partition size with an IEC suffix, for example 512MiB or 50GiB.",
				},
			},
		},
	}
}

// expandStorageLayout converts the Terraform model into the API request
// payload. A null mount is omitted, an explicitly empty one is sent as "".
func expandStorageLayout(ctx context.Context, model *StorageLayoutModel) *v1.StorageLayout {
	if model == nil {
		return nil
	}

	layout := &v1.StorageLayout{
		Disks: make([]*v1.StorageDisk, 0, len(model.Disks)),
	}
	for _, disk := range model.Disks {
		layout.Disks = append(layout.Disks, &v1.StorageDisk{
			Id:         disk.Id.ValueString(),
			Type:       strings.ToUpper(disk.Type.ValueString()),
			SizeGib:    disk.Size.ValueInt64(),
			Partitions: expandPartitions(disk.Partitions),
		})
	}
	if len(model.Raids) > 0 {
		layout.Raids = make([]*v1.StorageRaid, 0, len(model.Raids))
	}
	for _, raid := range model.Raids {
		var members []string
		raid.Members.ElementsAs(ctx, &members, false)
		layout.Raids = append(layout.Raids, &v1.StorageRaid{
			Id:         raid.Id.ValueString(),
			Type:       strings.ToUpper(raid.Type.ValueString()),
			Members:    members,
			Partitions: expandPartitions(raid.Partitions),
		})
	}
	return layout
}

func expandPartitions(models []StoragePartitionModel) []*v1.StoragePartition {
	if len(models) == 0 {
		return nil
	}
	partitions := make([]*v1.StoragePartition, 0, len(models))
	for _, p := range models {
		var mount *string
		if !p.Mount.IsNull() {
			value := p.Mount.ValueString()
			mount = &value
		}
		partitions = append(partitions, &v1.StoragePartition{
			Mount:  mount,
			Fstype: strings.ToUpper(p.Fs.ValueString()),
			Size:   p.Size.ValueString(),
		})
	}
	return partitions
}

// flattenStorageLayout converts the API response into the Terraform model.
// GetServer returns the user-submitted document, so a null layout stays null.
func flattenStorageLayout(ctx context.Context, layout *v1.StorageLayout) *StorageLayoutModel {
	if layout == nil {
		return nil
	}

	model := &StorageLayoutModel{
		Disks: make([]StorageDiskModel, 0, len(layout.Disks)),
	}
	for _, disk := range layout.Disks {
		if disk == nil {
			continue
		}
		model.Disks = append(model.Disks, StorageDiskModel{
			Id:         types.StringValue(disk.Id),
			Type:       types.StringValue(strings.ToLower(disk.Type)),
			Size:       types.Int64Value(disk.SizeGib),
			Partitions: flattenPartitions(disk.Partitions),
		})
	}
	for _, raid := range layout.Raids {
		if raid == nil {
			continue
		}
		members, diags := types.ListValueFrom(ctx, types.StringType, raid.Members)
		_ = diags // string conversion cannot fail
		model.Raids = append(model.Raids, StorageRaidModel{
			Id:         types.StringValue(raid.Id),
			Type:       types.StringValue(strings.ToLower(raid.Type)),
			Members:    members,
			Partitions: flattenPartitions(raid.Partitions),
		})
	}
	return model
}

func flattenPartitions(partitions []*v1.StoragePartition) []StoragePartitionModel {
	if len(partitions) == 0 {
		return nil
	}
	models := make([]StoragePartitionModel, 0, len(partitions))
	for _, p := range partitions {
		if p == nil {
			continue
		}
		mount := types.StringNull()
		if p.Mount != nil {
			mount = types.StringValue(*p.Mount)
		}
		models = append(models, StoragePartitionModel{
			Mount: mount,
			Fs:    types.StringValue(strings.ToLower(p.Fstype)),
			Size:  types.StringValue(p.Size),
		})
	}
	return models
}

// equalStorageLayout compares two layout models for the reprovision trigger.
// Disks and raids are unordered (matched by id); partition order is
// significant — it determines placement on the device.
func equalStorageLayout(a, b *StorageLayoutModel) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}

	disks := slices.Clone(a.Disks)
	otherDisks := slices.Clone(b.Disks)
	slices.SortFunc(disks, func(x, y StorageDiskModel) int {
		return cmp.Compare(x.Id.ValueString(), y.Id.ValueString())
	})
	slices.SortFunc(otherDisks, func(x, y StorageDiskModel) int {
		return cmp.Compare(x.Id.ValueString(), y.Id.ValueString())
	})
	if !equalDisks(disks, otherDisks) {
		return false
	}

	raids := slices.Clone(a.Raids)
	otherRaids := slices.Clone(b.Raids)
	slices.SortFunc(raids, func(x, y StorageRaidModel) int {
		return cmp.Compare(x.Id.ValueString(), y.Id.ValueString())
	})
	slices.SortFunc(otherRaids, func(x, y StorageRaidModel) int {
		return cmp.Compare(x.Id.ValueString(), y.Id.ValueString())
	})
	return equalRaids(raids, otherRaids)
}

func equalDisks(a, b []StorageDiskModel) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Id.Equal(b[i].Id) ||
			!strings.EqualFold(a[i].Type.ValueString(), b[i].Type.ValueString()) ||
			!a[i].Size.Equal(b[i].Size) ||
			!equalPartitions(a[i].Partitions, b[i].Partitions) {
			return false
		}
	}
	return true
}

func equalRaids(a, b []StorageRaidModel) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Id.Equal(b[i].Id) ||
			!strings.EqualFold(a[i].Type.ValueString(), b[i].Type.ValueString()) ||
			!a[i].Members.Equal(b[i].Members) ||
			!equalPartitions(a[i].Partitions, b[i].Partitions) {
			return false
		}
	}
	return true
}

func equalPartitions(a, b []StoragePartitionModel) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Mount.Equal(b[i].Mount) ||
			!strings.EqualFold(a[i].Fs.ValueString(), b[i].Fs.ValueString()) ||
			!a[i].Size.Equal(b[i].Size) {
			return false
		}
	}
	return true
}
