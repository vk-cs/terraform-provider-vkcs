package baremetal

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/baremetal/v1/flavors"
)

var (
	_ datasource.DataSource              = &FlavorsDataSource{}
	_ datasource.DataSourceWithConfigure = &FlavorsDataSource{}
)

func NewFlavorsDataSource() datasource.DataSource {
	return &FlavorsDataSource{}
}

type FlavorsDataSource struct {
	config clients.Config
}

type FlavorsDataSourceModel struct {
	Region  types.String  `tfsdk:"region"`
	Flavors []FlavorModel `tfsdk:"flavors"`
}

type FlavorModel struct {
	Name            types.String `tfsdk:"name"`
	CpuModel        types.String `tfsdk:"cpu_model"`
	CpuCores        types.Int64  `tfsdk:"cpu_cores"`
	RamSize         types.Int64  `tfsdk:"ram_size"`
	SsdSize         types.Int64  `tfsdk:"ssd_size"`
	HddSize         types.Int64  `tfsdk:"hdd_size"`
	BondVlanCapable types.Bool   `tfsdk:"bond_vlan_capable"`
}

func (d *FlavorsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_baremetal_flavors"
}

func (d *FlavorsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "The region to fetch the bare metal flavor from, defaults to the provider's region.",
			},
			"flavors": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the flavor.",
						},
						"cpu_model": schema.StringAttribute{
							Computed:    true,
							Description: "The CPU model.",
						},
						"cpu_cores": schema.Int64Attribute{
							Computed:    true,
							Description: "CPU core count including hyper-threading.",
						},
						"ram_size": schema.Int64Attribute{
							Computed:    true,
							Description: "RAM in gigabytes.",
						},
						"ssd_size": schema.Int64Attribute{
							Computed:    true,
							Description: "SSD size in gigabytes.",
						},
						"hdd_size": schema.Int64Attribute{
							Computed:    true,
							Description: "HDD size in gigabytes.",
						},
						"bond_vlan_capable": schema.BoolAttribute{
							Computed:    true,
							Description: "Bond and VLAN capable.",
						},
					},
				},
				Description: "Available Baremetal Flavors.",
			},
		},
		Description: "Use this data source to get a list of available VKCS Baremetal Flavors.",
	}
}

func (d *FlavorsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.config = req.ProviderData.(clients.Config)
}

func (d *FlavorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FlavorsDataSourceModel

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

	tflog.Debug(ctx, "Calling baremetal API to get list of flavors")

	flavorPages, err := flavors.List(client, nil).AllPages()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Baremetal API", err.Error())
		return
	}

	fvs, err := flavors.ExtractFlavors(flavorPages)
	if err != nil {
		resp.Diagnostics.AddError("Error extract flavors", err.Error())
		return
	}

	data.Region = types.StringValue(region)
	data.Flavors = flattenFlavors(fvs)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenFlavors(items []flavors.Flavor) (r []FlavorModel) {
	for _, item := range items {
		elem := FlavorModel{
			Name:            types.StringValue(item.FlavorName),
			CpuModel:        types.StringValue(item.CpuModel),
			CpuCores:        types.Int64Value(item.CpuCores),
			RamSize:         types.Int64Value(item.RamGb),
			BondVlanCapable: types.BoolValue(item.BondAndVlanCapable),
		}
		for _, disk := range item.Disks {
			if disk.Type == flavors.DiskTypeSSD {
				elem.SsdSize = types.Int64Value(disk.Size)
			}
			if disk.Type == flavors.DiskTypeHDD {
				elem.HddSize = types.Int64Value(disk.Size)
			}
		}
		r = append(r, elem)
	}

	sort.SliceStable(r, func(i, j int) bool {
		return r[i].Name.ValueString() < r[j].Name.ValueString()
	})

	return
}
