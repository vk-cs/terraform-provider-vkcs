package kubernetes_v2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
	dskubengv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/v2/datasource_kubernetes_node_group_v2"
)

var (
	_ datasource.DataSource              = (*kubernetesNodeGroupV2DataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*kubernetesNodeGroupV2DataSource)(nil)
)

func NewKubernetesNodeGroupV2DataSource() datasource.DataSource {
	return &kubernetesNodeGroupV2DataSource{}
}

type kubernetesNodeGroupV2DataSource struct {
	config clients.Config
}

func (d *kubernetesNodeGroupV2DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_node_group_v2"
}

func (d *kubernetesNodeGroupV2DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dskubengv2.KubernetesNodeGroupV2DataSourceSchema(ctx)
}

func (d *kubernetesNodeGroupV2DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *kubernetesNodeGroupV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dskubengv2.KubernetesNodeGroupV2Model

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
		resp.Diagnostics.AddError("Error creating Container Infra V2 API client for datasource vkcs_kubernetes_node_group_v2", err.Error())
		return
	}

	nodeGroupID := data.Id.ValueString()
	if nodeGroupID == "" {
		resp.Diagnostics.AddError("Error finding a node group", "id must be specified")
		return
	}

	tflog.Trace(ctx, "Calling Container Infra V2 API to get node group", map[string]interface{}{"node_group_id": fmt.Sprintf("%#v", nodeGroupID)})

	nodeGroup, err := nodegroups.Get(client, nodeGroupID).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling Container Infra V2 API to get node group", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Container Infra V2 API to get node group", map[string]interface{}{
		"node_group_id": fmt.Sprintf("%#v", nodeGroupID),
		"node_group":    fmt.Sprintf("%#v", nodeGroup),
	})

	// Use converter function to flatten node group data
	data, diags := dskubengv2.FlattenNodeGroup(ctx, nodeGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Region = types.StringValue(region)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
