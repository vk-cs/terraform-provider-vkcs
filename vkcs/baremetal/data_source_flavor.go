package baremetal

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/baremetal/v1/flavors"
)

var _ datasource.DataSource = &FlavorDataSource{}
var _ datasource.DataSourceWithConfigure = &FlavorDataSource{}

func NewFlavorDataSource() datasource.DataSource {
	return &FlavorDataSource{}
}

type FlavorDataSource struct {
	config clients.Config
}

type FlavorDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Region          types.String `tfsdk:"region"`
	Name            types.String `tfsdk:"name"`
	CpuModel        types.String `tfsdk:"cpu_model"`
	CpuCores        types.Int64  `tfsdk:"cpu_cores"`
	RamSize         types.Int64  `tfsdk:"ram_size"`
	SsdSize         types.Int64  `tfsdk:"ssd_size"`
	HddSize         types.Int64  `tfsdk:"hdd_size"`
	BondVlanCapable types.Bool   `tfsdk:"bond_vlan_capable"`
}

func (d *FlavorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_baremetal_flavor"
}

func (d *FlavorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The UUID of the flavor.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("name"),
						path.MatchRoot("cpu_model"),
						path.MatchRoot("cpu_cores"),
						path.MatchRoot("ram_size"),
						path.MatchRoot("ssd_size"),
						path.MatchRoot("hdd_size"),
						path.MatchRoot("bond_vlan_capable"),
					}...),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the flavor.",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "The region to fetch the bare metal flavor from, defaults to the provider's region.",
			},
			"cpu_model": schema.StringAttribute{
				Optional:    true,
				Description: "The CPU model.",
			},
			"cpu_cores": schema.Int64Attribute{
				Optional:    true,
				Description: "CPU core count including hyper-threading.",
			},
			"ram_size": schema.Int64Attribute{
				Optional:    true,
				Description: "RAM in gigabytes.",
			},
			"ssd_size": schema.Int64Attribute{
				Optional:    true,
				Description: "SSD size in gigabytes.",
			},
			"hdd_size": schema.Int64Attribute{
				Optional:    true,
				Description: "HDD size in gigabytes.",
			},
			"bond_vlan_capable": schema.BoolAttribute{
				Optional:    true,
				Description: "Bond and VLAN capable.",
			},
		},
		Description: "Use this data source to get information about a VKCS Baremetal Flavor.",
	}
}

func (d *FlavorDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.config = req.ProviderData.(clients.Config)
}

func (d *FlavorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FlavorDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	client, err := d.config.BareMetalV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS Baremetal API client", err.Error())
		return
	}

	flavor, err := getFlavor(client, data)
	if err != nil {
		resp.Diagnostics.AddError("Error getting flavor", err.Error())
		return
	}

	data.ID = types.StringValue(flavor.FlavorId)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(flavor.FlavorName)
	data.CpuModel = types.StringValue(flavor.CpuModel)
	data.CpuCores = types.Int64Value(flavor.CpuCores)
	data.RamSize = types.Int64Value(flavor.RamGb)
	for _, disk := range flavor.Disks {
		if disk.Type == flavors.DiskTypeSSD {
			data.SsdSize = types.Int64Value(disk.Size)
		}
		if disk.Type == flavors.DiskTypeHDD {
			data.HddSize = types.Int64Value(disk.Size)
		}
	}
	data.BondVlanCapable = types.BoolValue(flavor.BondAndVlanCapable)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getFlavor(client *gophercloud.ServiceClient, data FlavorDataSourceModel) (*flavors.Flavor, error) {
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		return flavors.Get(client, data.ID.ValueString()).Extract()
	}

	flavorPages, err := flavors.List(client, &flavors.ListOpts{
		CpuModel:    data.CpuModel.ValueStringPointer(),
		CpuCoresMin: data.CpuCores.ValueInt64Pointer(),
		CpuCoresMax: data.CpuCores.ValueInt64Pointer(),
		RamSizeMin:  data.RamSize.ValueInt64Pointer(),
		RamSizeMax:  data.RamSize.ValueInt64Pointer(),
		SsdSizeMin:  data.SsdSize.ValueInt64Pointer(),
		SsdSizeMax:  data.SsdSize.ValueInt64Pointer(),
		HddSizeMin:  data.HddSize.ValueInt64Pointer(),
		HddSizeMax:  data.HddSize.ValueInt64Pointer(),
	}).AllPages()
	if err != nil {
		return nil, err
	}

	fvs, err := flavors.ExtractFlavors(flavorPages)
	if err != nil {
		return nil, err
	}

	filterFvs := filterFlavors(fvs, data)

	if len(filterFvs) == 0 {
		return nil, fmt.Errorf("no flavors found by params")
	}

	if len(filterFvs) > 1 {
		return nil, fmt.Errorf("multiple flavors found by params")
	}

	return &filterFvs[0], nil
}

func filterFlavors(items []flavors.Flavor, data FlavorDataSourceModel) (r []flavors.Flavor) {
	for _, item := range items {
		if flavorMatches(data, item) {
			r = append(r, item)
		}
	}

	return
}

func flavorMatches(data FlavorDataSourceModel, item flavors.Flavor) bool {
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		if item.FlavorName != data.Name.ValueString() {
			return false
		}
	}

	if !data.BondVlanCapable.IsNull() && !data.BondVlanCapable.IsUnknown() {
		if item.BondAndVlanCapable != data.BondVlanCapable.ValueBool() {
			return false
		}
	}

	return true
}
