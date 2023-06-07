package kubernetes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	v1 "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfraaddons/v1"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfraaddons/v1/addons"
)

var _ datasource.DataSource = &AddonDataSource{}
var _ datasource.DataSourceWithConfigure = &AddonDataSource{}

func NewAddonDatasource() datasource.DataSource {
	return &AddonDataSource{}
}

type AddonDataSource struct {
	config clients.Config
}

type AddonDataSourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Region              types.String `tfsdk:"region"`
	ClusterID           types.String `tfsdk:"cluster_id"`
	Name                types.String `tfsdk:"name"`
	Version             types.String `tfsdk:"version"`
	ConfigurationValues types.String `tfsdk:"configuration_values"`
	Installed           types.Bool   `tfsdk:"installed"`
}

func (d *AddonDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_kubernetes_addon"
}

func (d *AddonDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

			"name": schema.StringAttribute{
				Required:    true,
				Description: "An addon name to filter by.",
			},

			"version": schema.StringAttribute{
				Required:    true,
				Description: "An addon version to filter by.",
			},

			"configuration_values": schema.StringAttribute{
				Computed: true,
				Description: "Configuration code for the addon. If the addon was installed in the cluster, " +
					"this value is the user-provided configuration code, otherwise it is a template for this cluster.",
			},

			"installed": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the addon was installed in the cluster.",
			},
		},

		Description: "Provides a kubernetes cluster addon datasource. This can be used to get information about an VKCS cluster addon.",
	}
}

func (d *AddonDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *AddonDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AddonDataSourceModel

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

	var allAddons []v1.Addon
	allAddons = append(allAddons, installedAddons...)
	allAddons = append(allAddons, availableAddons...)

	name, version := data.Name.ValueString(), data.Version.ValueString()

	tflog.Debug(ctx, "Filtering retrieved addons", map[string]interface{}{"name": name, version: version})
	filteredAddons := filterAddons(allAddons, name, version)
	tflog.Debug(ctx, "Filtered retrieved addons", map[string]interface{}{"filtered_addons": fmt.Sprintf("%#v", filteredAddons)})

	if len(filteredAddons) < 1 {
		resp.Diagnostics.AddError("Error filtering addons", "Your query returned no results."+
			"Please change your search criteria and try again.")
		return
	}

	if len(filteredAddons) > 1 {
		resp.Diagnostics.AddError("Error filtering addons",
			"Your query returned more than one result. Please try a more specific search criteria.")
		return
	}

	addon := filteredAddons[0]

	var (
		cfgValues   string
		isInstalled bool
	)

	if addonIn(addon, availableAddons) {
		availableAddon, err := addons.GetAvailableAddon(client, clusterID, addon.ID).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error getting available addon", err.Error())
			return
		}
		addon.ValuesTemplate = availableAddon.ValuesTemplate
		cfgValues = addon.ValuesTemplate
	}

	if addonIn(addon, installedAddons) {
		cfgValues, err = getConfigurationValues(addon.ID, clusterAddons)
		if err != nil {
			resp.Diagnostics.AddError("Error getting configuration values for addon", err.Error())
			return
		}
		isInstalled = true
	}

	data.ID = types.StringValue(addon.ID)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(addon.Name)
	data.Version = types.StringValue(addon.ChartVersion)
	data.ConfigurationValues = types.StringValue(cfgValues)
	data.Installed = types.BoolValue(isInstalled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
