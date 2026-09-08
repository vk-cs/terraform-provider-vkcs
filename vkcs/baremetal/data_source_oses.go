package baremetal

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/baremetal/v1/images"
)

var (
	_ datasource.DataSource              = &OSesDataSource{}
	_ datasource.DataSourceWithConfigure = &OSesDataSource{}
)

func NewOSesDataSource() datasource.DataSource {
	return &OSesDataSource{}
}

type OSesDataSource struct {
	config clients.Config
}

type OSesDataSourceModel struct {
	Region types.String `tfsdk:"region"`
	OSes   []OSModel    `tfsdk:"oses"`
}

type OSModel struct {
	Name     types.String `tfsdk:"name"`
	Version  types.String `tfsdk:"version"`
	RaidType types.String `tfsdk:"raid_type"`
}

func (d *OSesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_baremetal_oses"
}

func (d *OSesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "The region to fetch the bare metal OSes from, defaults to the provider's region.",
			},
			"oses": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the OS.",
						},
						"version": schema.StringAttribute{
							Computed:    true,
							Description: "The version of the OS.",
						},
						"raid_type": schema.StringAttribute{
							Computed:    true,
							Description: "The raid type of the OS.",
						},
					},
				},
				Description: "Available Baremetal OSes.",
			},
		},
		Description: "Use this data source to get a list of available VKCS Baremetal OSes.",
	}
}

func (d *OSesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.config = req.ProviderData.(clients.Config)
}

func (d *OSesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OSesDataSourceModel

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

	tflog.Debug(ctx, "Calling baremetal API to get list of images")

	imagePages, err := images.List(client).AllPages()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Baremetal API", err.Error())
		return
	}

	imgs, err := images.ExtractImages(imagePages)
	if err != nil {
		resp.Diagnostics.AddError("Error extract images", err.Error())
		return
	}

	data.Region = types.StringValue(region)
	data.OSes = flattenImages(imgs)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenImages(items []images.Image) (r []OSModel) {
	for _, item := range items {
		r = append(r, OSModel{
			Name:     types.StringValue(item.OsType),
			Version:  types.StringValue(item.OsVersion),
			RaidType: types.StringValue(item.RaidType),
		})
	}

	sort.SliceStable(r, func(i, j int) bool {
		if r[i].Name.ValueString() != r[j].Name.ValueString() {
			return r[i].Name.ValueString() < r[j].Name.ValueString()
		}

		return r[i].Version.ValueString() < r[j].Version.ValueString()
	})

	return
}
