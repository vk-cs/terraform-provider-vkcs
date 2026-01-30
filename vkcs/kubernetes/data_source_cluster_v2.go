package kubernetes

import (
	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	dskubeclusterv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/datasource_kubernetes_cluster_v2"
)

var (
	_ datasource.DataSource = (*kubernetesClusterV2DataSource)(nil)
)

func NewKubernetesClusterV2DataSource() datasource.DataSource {
	return &kubernetesClusterV2DataSource{}
}

type kubernetesClusterV2DataSource struct {
	config clients.Config
}

func (d *kubernetesClusterV2DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_cluster_v2"
}

func (d *kubernetesClusterV2DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dskubeclusterv2.KubernetesClusterV2DataSourceSchema(ctx)
}

func (d *kubernetesClusterV2DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *kubernetesClusterV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dskubeclusterv2.KubernetesClusterV2Model

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
		resp.Diagnostics.AddError("Error creating API client for datasource vkcs_kubernetes_cluster_v2", err.Error())
		return
	}

	// Read the cluster from Managed K8S API to populate fields
	apiCluster, kubeconfig, diags := d.readCluster(client, data.Id.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.UpdateFromCluster(ctx, apiCluster, kubeconfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesClusterV2DataSource) readCluster(client *gophercloud.ServiceClient, clusterID string) (cluster *clusters.Cluster, clusterKubeconfig *string, diags diag.Diagnostics) {
	// Get cluster data
	cluster, err := clusters.Get(client, clusterID).Extract()
	if err != nil {
		diags.AddError("Error reading Kubernetes cluster", err.Error())
		return
	}

	// Get cluster kubeconfig
	kubeconfig, err := clusters.GetKubeconfig(client, clusterID)
	if err != nil {
		diags.AddError("Error retrieving cluster kubeconfig", err.Error())
		return
	}

	clusterKubeconfig = &kubeconfig
	return
}
