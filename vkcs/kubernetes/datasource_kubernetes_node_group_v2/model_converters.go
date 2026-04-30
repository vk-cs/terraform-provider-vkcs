package datasource_kubernetes_node_group_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
	mshared "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/models_shared"
)

func (m *KubernetesNodeGroupV2Model) UpdateFromNodeGroup(ctx context.Context, nodeGroup *nodegroups.NodeGroup) (diags diag.Diagnostics) {
	// Set basic fields
	m.Id = types.StringValue(nodeGroup.ID)
	m.Uuid = types.StringValue(nodeGroup.UUID)
	m.ClusterId = types.StringValue(nodeGroup.ClusterID)
	m.CreatedAt = types.StringValue(nodeGroup.CreatedAt)
	m.Name = types.StringValue(nodeGroup.Name)

	// Set availability zone
	m.AvailabilityZone, diags = FlattenAvailabilityZone(nodeGroup.Zones)
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
	m.ScaleType, m.FixedScaleNodeCount, m.AutoScaleNodeCount, m.AutoScaleMinSize, m.AutoScaleMaxSize, diags = FlattenScaleSpec(nodeGroup.ScaleSpec)
	if diags.HasError() {
		return diags
	}

	// Set taints
	m.Taints, diags = FlattenTaints(ctx, nodeGroup.Taints)
	if diags.HasError() {
		return diags
	}

	// Set root disk info
	m.DiskType, m.DiskSize, diags = FlattenDiskType(nodeGroup.DiskType)
	if diags.HasError() {
		return diags
	}

	return
}

func FlattenTaints(ctx context.Context, taints []nodegroups.Taint) (types.Set, diag.Diagnostics) {
	if len(taints) == 0 {
		return types.SetNull(TaintsValue{}.Type(ctx)), nil
	}

	resList := make([]attr.Value, len(taints))
	for i, taint := range taints {
		resList[i] = TaintsValue{
			Key:    types.StringValue(taint.Key),
			Value:  types.StringValue(taint.Value),
			Effect: types.StringValue(taint.Effect),
			state:  attr.ValueStateKnown,
		}
	}

	return types.SetValue(TaintsValue{}.Type(ctx), resList)
}

func FlattenScaleSpec(scaleSpec nodegroups.ScaleSpec) (
	scaleType types.String,
	fixedScaleNodeCount types.Int64,
	autoScaleNodeCount types.Int64,
	autoScaleMinSize types.Int64,
	autoScaleMaxSize types.Int64,
	diags diag.Diagnostics,
) {
	switch {
	case scaleSpec.FixedScale != nil:
		scaleType = types.StringValue("fixed_scale")
		fixedScaleNodeCount = types.Int64Value(int64(scaleSpec.FixedScale.Size))
		autoScaleNodeCount = types.Int64Null()
		autoScaleMinSize = types.Int64Null()
		autoScaleMaxSize = types.Int64Null()
	case scaleSpec.AutoScale != nil:
		scaleType = types.StringValue("auto_scale")
		autoScaleNodeCount = types.Int64Value(int64(scaleSpec.AutoScale.Size))
		autoScaleMinSize = types.Int64Value(int64(scaleSpec.AutoScale.MinSize))
		autoScaleMaxSize = types.Int64Value(int64(scaleSpec.AutoScale.MaxSize))
		fixedScaleNodeCount = types.Int64Null()
	default:
		scaleType = types.StringNull()
		fixedScaleNodeCount = types.Int64Null()
		autoScaleNodeCount = types.Int64Null()
		autoScaleMinSize = types.Int64Null()
		autoScaleMaxSize = types.Int64Null()
	}

	return
}

func FlattenAvailabilityZone(zones []string) (types.String, diag.Diagnostics) {
	if len(zones) == 0 {
		return types.StringNull(), nil
	}
	return types.StringValue(zones[0]), nil
}

func FlattenDiskType(diskType nodegroups.DiskType) (types.String, types.Int64, diag.Diagnostics) {
	if diskType.CinderVolumeType.Type != "" && diskType.CinderVolumeType.Size != 0 {
		return types.StringValue(diskType.CinderVolumeType.Type), types.Int64Value(int64(diskType.CinderVolumeType.Size)), nil
	}
	return types.StringNull(), types.Int64Null(), nil
}
