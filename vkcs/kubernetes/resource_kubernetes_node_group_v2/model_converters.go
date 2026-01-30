package resource_kubernetes_node_group_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
	dskubengv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/datasource_kubernetes_node_group_v2"
	mshared "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/models_shared"
)

func (m *KubernetesNodeGroupV2Model) UpdateFromNodeGroup(ctx context.Context, nodeGroup *nodegroups.NodeGroup) (diags diag.Diagnostics) {
	if nodeGroup == nil {
		return nil
	}

	// Set basic fields
	m.Id = types.StringValue(nodeGroup.ID)
	m.Uuid = types.StringValue(nodeGroup.UUID)
	m.ClusterId = types.StringValue(nodeGroup.ClusterID)
	m.CreatedAt = types.StringValue(nodeGroup.CreatedAt)
	m.Name = types.StringValue(nodeGroup.Name)

	// Set availability zone
	m.AvailabilityZone, diags = dskubengv2.FlattenAvailabilityZone(nodeGroup.Zones)
	if diags.HasError() {
		return diags
	}

	// Set labels
	m.Labels, diags = mshared.FlattenStringMap(nodeGroup.Labels)
	if diags.HasError() {
		return diags
	}

	// Set parallel upgrade chunk
	m.ParallelUpgradeChunk = types.Int64Value(int64(nodeGroup.ParallelUpgradeChunk))

	// Set node flavor
	m.NodeFlavor = types.StringValue(nodeGroup.VMEngine.NovaEngine.FlavorID)

	// Set scale type and related fields
	m.ScaleType, m.FixedScaleNodeCount, m.AutoScaleNodeCount, m.AutoScaleMinSize, m.AutoScaleMaxSize, diags = dskubengv2.FlattenScaleSpec(nodeGroup.ScaleSpec)
	if diags.HasError() {
		return diags
	}

	// Set taints
	m.Taints, diags = dskubengv2.FlattenTaints(ctx, nodeGroup.Taints)
	if diags.HasError() {
		return diags
	}

	// Set root disk info
	m.DiskType, m.DiskSize, diags = dskubengv2.FlattenDiskType(nodeGroup.DiskType)
	if diags.HasError() {
		return diags
	}

	return
}
