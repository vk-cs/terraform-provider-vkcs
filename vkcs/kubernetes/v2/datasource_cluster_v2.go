package kubernetes_v2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	dskubeclusterv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/v2/datasource_kubernetes_cluster_v2"
)

var (
	_ datasource.DataSource              = (*kubernetesClusterV2DataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*kubernetesClusterV2DataSource)(nil)
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

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	client, err := d.config.ContainerInfraV2Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Container Infra V2 API client for datasource vkcs_kubernetes_cluster_v2", err.Error())
		return
	}

	clusterID := data.Id.ValueString()
	if clusterID == "" {
		resp.Diagnostics.AddError("Error finding a cluster", "Cluster ID must be specified")
		return
	}

	tflog.Trace(ctx, "Calling Container Infra V2 API to get cluster", map[string]interface{}{"cluster_id": fmt.Sprintf("%#v", clusterID)})

	cluster, err := clusters.Get(client, clusterID).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling Container Infra V2 API to get cluster", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Container Infra V2 API to get cluster", map[string]interface{}{
		"cluster_id": fmt.Sprintf("%#v", clusterID),
		"cluster":    fmt.Sprintf("%#v", cluster),
	})

	// Use converter function to flatten cluster data
	data, diags := dskubeclusterv2.FlattenCluster(ctx, cluster)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Region = types.StringValue(region)

	// Get cluster kubeconfig
	kubeconfig, err := clusters.GetKubeconfig(client, clusterID)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving cluster kubeconfig", err.Error())
		return
	}
	data.K8sConfig = types.StringValue(kubeconfig)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
