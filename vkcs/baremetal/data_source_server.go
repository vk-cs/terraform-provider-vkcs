package baremetal

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/baremetal/v1/servers"
)

var _ datasource.DataSource = &ServerDataSource{}
var _ datasource.DataSourceWithConfigure = &ServerDataSource{}

func NewServerDataSource() datasource.DataSource {
	return &ServerDataSource{}
}

type ServerDataSource struct {
	config clients.Config
}

type ServerDataSourceModel struct {
	ID               types.String           `tfsdk:"id"`
	Region           types.String           `tfsdk:"region"`
	Name             types.String           `tfsdk:"name"`
	AvailabilityZone types.String           `tfsdk:"availability_zone"`
	CpuCores         types.List             `tfsdk:"cpu_cores"`
	CpuTypes         types.List             `tfsdk:"cpu_types"`
	IsLocked         types.Bool             `tfsdk:"is_locked"`
	LocalDiskSizes   types.List             `tfsdk:"local_disk_sizes"`
	RamMegabytes     types.Int64            `tfsdk:"ram_megabytes"`
	PowerState       types.String           `tfsdk:"power_state"`
	Tags             types.List             `tfsdk:"tags"`
	ImageID          types.String           `tfsdk:"image_id"`
	ImageName        types.String           `tfsdk:"image_name"`
	OsType           types.String           `tfsdk:"os_type"`
	RaidType         types.String           `tfsdk:"raid_type"`
	Status           types.String           `tfsdk:"status"`
	FlavorID         types.String           `tfsdk:"flavor_id"`
	TargetBootOrder  []TargetBootOrderModel `tfsdk:"target_boot_order"`
	LocalDisksInfo   []LocalDiskInfoModel   `tfsdk:"local_disks_info"`
}

type TargetBootOrderModel struct {
	BootDeviceType types.String `tfsdk:"boot_device_type"`
}

type LocalDiskInfoModel struct {
	Path  types.String `tfsdk:"path"`
	Size  types.Int64  `tfsdk:"size"`
	Type  types.String `tfsdk:"type"`
	Model types.String `tfsdk:"model"`
}

func (s *ServerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_baremetal_server"
}

func (s *ServerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the bare metal server.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the server.",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "The region to fetch the bare metal server from, defaults to the provider's region.",
			},
			"availability_zone": schema.StringAttribute{
				Computed:    true,
				Description: "The availability zone of this server.",
			},
			"cpu_cores": schema.ListAttribute{
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "CPU core count including hyper-threading.",
			},
			"cpu_types": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Server CPU type.",
			},
			"is_locked": schema.BoolAttribute{
				Computed: true,
				Description: "Shows whether the server is protected." +
					" The server cannot be deleted while this flag is set.",
			},
			"local_disk_sizes": schema.ListAttribute{
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "Local disk sizes in gigabytes.",
			},
			"ram_megabytes": schema.Int64Attribute{
				Computed:    true,
				Description: "Server memory size in megabytes.",
			},
			"power_state": schema.StringAttribute{
				Computed:    true,
				Description: "Server power state.",
			},
			"tags": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Server tags.",
			},
			"image_id": schema.StringAttribute{
				Computed:    true,
				Description: "The image ID used to create the server.",
			},
			"image_name": schema.StringAttribute{
				Computed:    true,
				Description: "The image name used to create the server.",
			},
			"os_type": schema.StringAttribute{
				Computed:    true,
				Description: "Server Operation System type.",
			},
			"raid_type": schema.StringAttribute{
				Computed:    true,
				Description: "Server raid type.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Server status.",
			},
			"flavor_id": schema.StringAttribute{
				Computed:    true,
				Description: "Flavor ID of the server.",
			},
			"target_boot_order": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Current server boot order.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"boot_device_type": schema.StringAttribute{
							Computed:    true,
							Description: "The boot device type.",
						},
					},
				},
			},
			"local_disks_info": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Information about server disks.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							Computed:    true,
							Description: "The path to the disk.",
						},
						"size": schema.Int64Attribute{
							Computed:    true,
							Description: "The size of the disk.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the disk.",
						},
						"model": schema.StringAttribute{
							Computed:    true,
							Description: "The model of the disk.",
						},
					},
				},
			},
		},
		Description: "Use this data source to get information about a VKCS bare metal server.",
	}
}

func (s *ServerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	s.config = req.ProviderData.(clients.Config)
}

func (s *ServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServerDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = s.config.GetRegion()
	}

	client, err := s.config.BareMetalV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS baremetal API client", err.Error())
		return
	}

	serverID := data.ID.ValueString()
	ctx = tflog.SetField(ctx, "server_id", serverID)

	tflog.Debug(ctx, "Calling VKCS baremetal API to retrieve server by id", map[string]interface{}{"id": serverID})
	server, err := servers.Get(client, serverID).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error reading baremetal server", err.Error())
		return
	}

	data.Name = types.StringValue(server.ServerName)
	data.Region = types.StringValue(region)
	data.AvailabilityZone = types.StringValue(server.AvailabilityZone)
	data.TargetBootOrder = flattenTargetBootOrder(server.TargetBootOrder)
	data.LocalDisksInfo = flattenLocalDiskInfo(server.LocalDisksInfo)
	data.IsLocked = types.BoolValue(server.IsLocked)
	data.RamMegabytes = types.Int64Value(server.MemoryMegabytes)
	data.PowerState = types.StringValue(server.PowerState)
	data.ImageID = types.StringPointerValue(server.ImageId)
	data.ImageName = types.StringPointerValue(server.ImageName)
	data.OsType = types.StringPointerValue(server.OsType)
	data.RaidType = types.StringPointerValue(server.RaidType)
	data.Status = types.StringValue(string(server.Status))
	data.FlavorID = types.StringPointerValue(server.FlavorId)

	var d diag.Diagnostics
	data.CpuCores, d = types.ListValueFrom(ctx, types.Int64Type, server.CpuCores)
	resp.Diagnostics.Append(d...)
	data.CpuTypes, d = types.ListValueFrom(ctx, types.StringType, server.CpuTypes)
	resp.Diagnostics.Append(d...)
	data.LocalDiskSizes, d = types.ListValueFrom(ctx, types.Int64Type, server.LocalDiskSizes)
	resp.Diagnostics.Append(d...)
	data.Tags, d = types.ListValueFrom(ctx, types.StringType, server.Tags)
	resp.Diagnostics.Append(d...)

	if d.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenTargetBootOrder(items []*servers.BootOrderListItem) []TargetBootOrderModel {
	r := make([]TargetBootOrderModel, len(items))
	for idx, item := range items {
		r[idx] = TargetBootOrderModel{
			BootDeviceType: types.StringValue(item.BootDeviceType),
		}
	}

	sort.SliceStable(r, func(i, j int) bool {
		return r[i].BootDeviceType.ValueString() > r[j].BootDeviceType.ValueString()
	})

	return r
}

func flattenLocalDiskInfo(info []*servers.LocalDiskInfo) []LocalDiskInfoModel {
	r := make([]LocalDiskInfoModel, len(info))
	for idx, item := range info {
		r[idx] = LocalDiskInfoModel{
			Path:  types.StringValue(item.Path),
			Size:  types.Int64Value(item.SizeGb),
			Type:  types.StringValue(item.Type),
			Model: types.StringPointerValue(item.Model),
		}
	}

	sort.SliceStable(r, func(i, j int) bool {
		return r[i].Path.ValueString() > r[j].Path.ValueString()
	})

	return r
}
