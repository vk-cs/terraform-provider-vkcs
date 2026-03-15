package resource_kubernetes_node_group_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
	dskubengv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/v2/datasource_kubernetes_node_group_v2"
)

// ToNodeGroupModel конвертирует nodegroups.NodeGroup в KubernetesNodeGroupV2Model
func ToNodeGroupModel(ctx context.Context, nodeGroup *nodegroups.NodeGroup) (KubernetesNodeGroupV2Model, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model KubernetesNodeGroupV2Model

	// Set basic fields
	model.Id = types.StringValue(nodeGroup.ID)
	model.Uuid = types.StringValue(nodeGroup.UUID)
	model.ClusterId = types.StringValue(nodeGroup.ClusterID)
	model.CreatedAt = types.StringValue(nodeGroup.CreatedAt)
	model.Name = types.StringValue(nodeGroup.Name)

	// Set availability zone
	model.AvailabilityZone, diags = dskubengv2.FlattenAvailabilityZone(nodeGroup.Zones)
	if diags.HasError() {
		return model, diags
	}

	// Set labels
	model.Labels, diags = dskubengv2.FlattenLabels(nodeGroup.Labels)
	if diags.HasError() {
		return model, diags
	}

	// Set parallel upgrade chunk
	model.ParallelUpgradeChunk = types.Int64Value(int64(nodeGroup.ParallelUpgradeChunk))

	// Set node flavor
	model.NodeFlavor = types.StringValue(nodeGroup.VMEngine.NovaEngine.FlavorID)

	// Set scale type and related fields
	model.ScaleType, model.FixedScaleNodeCount, model.AutoScaleNodeCount, model.AutoScaleMinSize, model.AutoScaleMaxSize, diags = dskubengv2.FlattenScaleSpec(nodeGroup.ScaleSpec)
	if diags.HasError() {
		return model, diags
	}

	// Set taints
	model.Taints, diags = dskubengv2.FlattenTaints(ctx, nodeGroup.Taints)
	if diags.HasError() {
		return model, diags
	}

	// Set root disk info
	model.DiskType, model.DiskSize, diags = dskubengv2.FlattenDiskType(nodeGroup.DiskType)
	if diags.HasError() {
		return model, diags
	}

	return model, diags
}
