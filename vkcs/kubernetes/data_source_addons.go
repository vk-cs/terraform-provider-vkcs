package kubernetes

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gophercloud/utils/terraform/hashcode"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	v1 "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfraaddons/v1"
)

var _ datasource.DataSource = &AddonsDataSource{}
var _ datasource.DataSourceWithConfigure = &AddonsDataSource{}

func NewAddonsDatasource() datasource.DataSource {
	return &AddonsDataSource{}
}

type AddonsDataSource struct {
	config clients.Config
}

type AddonsDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Region    types.String `tfsdk:"region"`
	ClusterID types.String `tfsdk:"cluster_id"`
	Addons    []AddonModel `tfsdk:"addons"`
}

type AddonModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Version   types.String `tfsdk:"version"`
	Installed types.Bool   `tfsdk:"installed"`
}

func (d *AddonsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_kubernetes_addons"
}

func (d *AddonsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource.",
			},

			"region": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The region in which to obtain the service client. If omitted, the `region` argument of the provider is used.",
			},

			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the kubernetes cluster.",
			},

			"addons": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of an addon.",
						},

						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of an addon.",
						},

						"version": schema.StringAttribute{
							Computed:    true,
							Description: "Version of an addon.",
						},

						"installed": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether an addon was installed in the cluster.",
						},
					},
				},
			},
		},

		Description: "Provides a kubernetes cluster addons datasource. This can be used to get information about an VKCS cluster addons.",
	}
}

func (d *AddonsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *AddonsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AddonsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	client, err := d.config.ContainerInfraAddonsV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Kubernetes Addons API client", err.Error())
		return
	}

	clusterID := data.ClusterID.ValueString()
	ctx = tflog.SetField(ctx, "cluster_id", clusterID)

	availableAddons, clusterAddons, err := readAllClusterAddons(ctx, client, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading cluster addons", err.Error())
		return
	}

	installedAddons, err := extractAddonsFromClusterAddons(clusterAddons)
	if err != nil {
		resp.Diagnostics.AddError("Error extracting addons from the cluster addons", err.Error())
		return
	}

	availableAddons = removeInstalledAddons(availableAddons, installedAddons)

	flattenedAddons, diags := flattenAddons(availableAddons, installedAddons)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := hashcode.String(fmt.Sprintf("%s/addons", clusterID))

	data.ID = types.StringValue(strconv.FormatInt(int64(id), 10))
	data.Region = types.StringValue(region)
	data.Addons = flattenedAddons

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenAddons(availableAddons, installedAddons []v1.Addon) (r []AddonModel, diags diag.Diagnostics) {
	for _, a := range availableAddons {
		r = append(r, AddonModel{
			ID:        types.StringValue(a.ID),
			Name:      types.StringValue(a.Name),
			Version:   types.StringValue(a.ChartVersion),
			Installed: types.BoolValue(false),
		})
	}
	for _, a := range installedAddons {
		r = append(r, AddonModel{
			ID:        types.StringValue(a.ID),
			Name:      types.StringValue(a.Name),
			Version:   types.StringValue(a.ChartVersion),
			Installed: types.BoolValue(true),
		})
	}
	return
}
