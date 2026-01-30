package kubernetes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	dskubeversionsv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/datasource_kubernetes_versions_v2"
)

var (
	_ datasource.DataSource = (*kubernetesClusterVersionV2DataSource)(nil)
)

func NewKubernetesClusterVersionV2DataSource() datasource.DataSource {
	return &kubernetesClusterVersionV2DataSource{}
}

type kubernetesClusterVersionV2DataSource struct {
	config clients.Config
}

func (d *kubernetesClusterVersionV2DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_versions_v2"
}

func (d *kubernetesClusterVersionV2DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dskubeversionsv2.KubernetesVersionsV2DataSourceSchema(ctx)
}

func (d *kubernetesClusterVersionV2DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *kubernetesClusterVersionV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dskubeversionsv2.KubernetesVersionsV2Model

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
		resp.Diagnostics.AddError("Error creating API client for datasource vkcs_kubernetes_versions_v2", err.Error())
		return
	}

	// API call
	clusterVersions, err := clusters.GetListK8SVersion(client).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling API to get Kubernetes versions", err.Error())
		return
	}

	resp.Diagnostics.Append(data.UpdateFromClusterVersion(clusterVersions)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
