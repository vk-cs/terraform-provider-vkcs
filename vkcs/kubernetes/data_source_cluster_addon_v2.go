package kubernetes

import (
	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	addons "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/addons"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	dskubecladdonv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/datasource_kubernetes_cluster_addon_v2"
)

var (
	_ datasource.DataSource              = (*kubernetesClusterAddonV2DataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*kubernetesClusterAddonV2DataSource)(nil)
)

func NewKubernetesClusterAddonV2DataSource() datasource.DataSource {
	return &kubernetesClusterAddonV2DataSource{}
}

type kubernetesClusterAddonV2DataSource struct {
	config clients.Config
}

func (d *kubernetesClusterAddonV2DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_cluster_addon_v2"
}

func (d *kubernetesClusterAddonV2DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dskubecladdonv2.KubernetesClusterAddonV2DataSourceSchema(ctx)
}

func (d *kubernetesClusterAddonV2DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *kubernetesClusterAddonV2DataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data dskubecladdonv2.KubernetesClusterAddonV2Model

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	idSet := !data.Id.IsNull()
	clusterIDSet := !data.ClusterId.IsNull()
	baseAddonNameSet := !data.BaseAddonName.IsNull()

	clusterIDNameSet := clusterIDSet && baseAddonNameSet

	variants := 0
	if idSet {
		variants++
	}
	if clusterIDNameSet {
		variants++
	}

	if variants > 1 {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Specify only one of: 'id' or ('cluster_id' and 'base_addon_name').",
		)
		return
	}

	if variants == 0 {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"You must specify one of: 'id' or both 'cluster_id' and 'base_addon_name'.",
		)
		return
	}

	if (clusterIDSet && !baseAddonNameSet) || (!clusterIDSet && baseAddonNameSet) {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Both 'cluster_id' and 'base_addon_name' must be specified together.",
		)
	}

	if util.IsKnownValue(data.Id) && data.Id.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Attribute 'id' must be not empty.",
		)
	}

	if util.IsKnownValue(data.ClusterId) && data.ClusterId.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Attribute 'cluster_id' must be not empty.",
		)
	}

	if util.IsKnownValue(data.BaseAddonName) && data.BaseAddonName.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Attribute 'base_addon_name' must be not empty.",
		)
	}
}

func (d *kubernetesClusterAddonV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dskubecladdonv2.KubernetesClusterAddonV2Model

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the region in which to obtain the Managed K8S client
	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	// Init Managed K8S client
	client, err := d.config.ManagedK8SClient(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating API client for datasource vkcs_kubernetes_cluster_addon_v2", err.Error())
		return
	}

	// API call
	clusterAddon, diags := d.getClusterAddon(client, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(data.UpdateFromAddon(ctx, clusterAddon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *kubernetesClusterAddonV2DataSource) getClusterAddon(client *gophercloud.ServiceClient, data *dskubecladdonv2.KubernetesClusterAddonV2Model) (clusterAddon *addons.ClusterAddon, diags diag.Diagnostics) {
	if !data.Id.IsNull() {
		apiAddon, err := addons.GetClusterAddon(client, data.Id.ValueString()).Extract()
		if err != nil {
			diags.AddError("Error calling API to get cluster addon by ID", err.Error())
			return
		}
		clusterAddon = &apiAddon
	} else {
		apiAddon, err := addons.GetClusterAddonByClusterAndName(client, data.ClusterId.ValueString(), data.BaseAddonName.ValueString()).Extract()
		if err != nil {
			diags.AddError("Error calling API to get cluster addon by cluster ID and base addon name", err.Error())
			return
		}
		clusterAddon = &apiAddon
	}

	return
}
