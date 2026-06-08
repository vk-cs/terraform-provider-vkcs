package kubernetes

import (
	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	dskubengv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/datasource_kubernetes_node_group_v2"
)

var (
	_ datasource.DataSource                   = (*kubernetesNodeGroupV2DataSource)(nil)
	_ datasource.DataSourceWithConfigure      = (*kubernetesNodeGroupV2DataSource)(nil)
	_ datasource.DataSourceWithValidateConfig = (*kubernetesNodeGroupV2DataSource)(nil)
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

func (d *kubernetesNodeGroupV2DataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data dskubengv2.KubernetesNodeGroupV2Model

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	idSet := !data.Id.IsNull()
	uuidSet := !data.Uuid.IsNull()
	nameSet := !data.Name.IsNull()
	clusterIDSet := !data.ClusterId.IsNull()

	nameClusterSet := nameSet && clusterIDSet

	variants := 0
	if idSet {
		variants++
	}
	if uuidSet {
		variants++
	}
	if nameClusterSet {
		variants++
	}

	if variants > 1 {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Specify only one of: 'id', 'uuid', or ('name' and 'cluster_id').",
		)
		return
	}

	if variants == 0 {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"You must specify one of: 'id', 'uuid', or both 'name' and 'cluster_id'.",
		)
		return
	}

	if (nameSet && !clusterIDSet) || (!nameSet && clusterIDSet) {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Both 'name' and 'cluster_id' must be specified together.",
		)
	}

	if util.IsKnownValue(data.Id) && data.Id.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Attribute 'id' must be not empty.",
		)
	}

	if util.IsKnownValue(data.Uuid) && data.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Attribute 'uuid' must be not empty.",
		)
	}

	if util.IsKnownValue(data.Name) && data.Name.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Attribute 'name' must be not empty.",
		)
	}

	if util.IsKnownValue(data.ClusterId) && data.ClusterId.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Attribute 'cluster_id' must be not empty.",
		)
	}
}

func (d *kubernetesNodeGroupV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dskubengv2.KubernetesNodeGroupV2Model

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
		resp.Diagnostics.AddError("Error creating API client for datasource vkcs_kubernetes_node_group_v2", err.Error())
		return
	}

	// API call
	nodeGroup, diags := d.readNodeGroup(client, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(data.UpdateFromNodeGroup(ctx, nodeGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesNodeGroupV2DataSource) readNodeGroup(client *gophercloud.ServiceClient, data *dskubengv2.KubernetesNodeGroupV2Model) (nodeGroup *nodegroups.NodeGroup, diags diag.Diagnostics) {
	var requestID string
	if !data.Id.IsNull() {
		requestID = data.Id.ValueString()
	}
	if !data.Uuid.IsNull() {
		requestID = data.Uuid.ValueString()
	}

	var err error
	if requestID != "" {
		// Lookup by ID
		nodeGroup, err = nodegroups.GetByID(client, requestID).Extract()
		if err != nil {
			diags.AddError("Error calling Managed K8S API to get node group by ID", err.Error())
			return
		}
	} else {
		// Lookup by name + cluster ID
		nodeGroup, err = nodegroups.GetByName(client, data.ClusterId.ValueString(), data.Name.ValueString()).Extract()
		if err != nil {
			diags.AddError("Error calling Managed K8S API to get node group by name", err.Error())
			return
		}
	}

	return
}
