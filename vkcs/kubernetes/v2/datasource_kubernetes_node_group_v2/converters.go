package datasource_kubernetes_node_group_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
)

// FlattenLabels конвертирует map[string]string в types.Map
func FlattenLabels(labels map[string]string) (types.Map, diag.Diagnostics) {
	if labels == nil {
		return types.MapNull(types.StringType), nil
	}

	labelsMap := make(map[string]attr.Value, len(labels))
	for k, v := range labels {
		labelsMap[k] = types.StringValue(v)
	}

	return types.MapValue(types.StringType, labelsMap)
}

// FlattenTaints конвертирует []nodegroups.Taint в types.Set
func FlattenTaints(ctx context.Context, taints []nodegroups.Taint) (types.Set, diag.Diagnostics) {
	if len(taints) == 0 {
		return types.SetNull(TaintsValue{}.Type(ctx)), nil
	}

	taintValues := make([]attr.Value, len(taints))
	for i, taint := range taints {
		taintValues[i] = TaintsValue{
			Key:    types.StringValue(taint.Key),
			Value:  types.StringValue(taint.Value),
			Effect: types.StringValue(taint.Effect),
			state:  attr.ValueStateKnown,
		}
	}

	return types.SetValue(TaintsValue{}.Type(ctx), taintValues)
}

// FlattenScaleSpec конвертирует nodegroups.ScaleSpec в поля scale_type и связанные поля
func FlattenScaleSpec(scaleSpec nodegroups.ScaleSpec) (types.String, types.Int64, types.Int64, types.Int64, types.Int64, diag.Diagnostics) {
	var scaleType types.String
	var fixedScaleNodeCount types.Int64
	var autoScaleNodeCount types.Int64
	var autoScaleMinSize types.Int64
	var autoScaleMaxSize types.Int64

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

	return scaleType, fixedScaleNodeCount, autoScaleNodeCount, autoScaleMinSize, autoScaleMaxSize, nil
}

// FlattenAvailabilityZone конвертирует []string в types.String (берет первую зону)
func FlattenAvailabilityZone(zones []string) (types.String, diag.Diagnostics) {
	if len(zones) == 0 {
		return types.StringNull(), nil
	}
	return types.StringValue(zones[0]), nil
}

// FlattenDiskType конвертирует nodegroups.DiskType в disk_type и disk_size
func FlattenDiskType(diskType nodegroups.DiskType) (types.String, types.Int64, diag.Diagnostics) {
	if diskType.CinderVolumeType.Type != "" && diskType.CinderVolumeType.Size != 0 {
		return types.StringValue(diskType.CinderVolumeType.Type), types.Int64Value(int64(diskType.CinderVolumeType.Size)), nil
	}
	return types.StringNull(), types.Int64Null(), nil
}

// FlattenNodeGroup конвертирует nodegroups.NodeGroup в KubernetesNodeGroupV2Model
func FlattenNodeGroup(ctx context.Context, nodeGroup *nodegroups.NodeGroup) (KubernetesNodeGroupV2Model, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model KubernetesNodeGroupV2Model

	// Set basic fields
	model.Id = types.StringValue(nodeGroup.ID)
	model.Uuid = types.StringValue(nodeGroup.UUID)
	model.ClusterId = types.StringValue(nodeGroup.ClusterID)
	model.CreatedAt = types.StringValue(nodeGroup.CreatedAt)
	model.Name = types.StringValue(nodeGroup.Name)

	// Set availability zone
	model.AvailabilityZone, diags = FlattenAvailabilityZone(nodeGroup.Zones)
	if diags.HasError() {
		return model, diags
	}

	// Set labels
	model.Labels, diags = FlattenLabels(nodeGroup.Labels)
	if diags.HasError() {
		return model, diags
	}

	// Set parallel upgrade chunk
	model.ParallelUpgradeChunk = types.Int64Value(int64(nodeGroup.ParallelUpgradeChunk))

	// Set node flavor
	model.NodeFlavor = types.StringValue(nodeGroup.VMEngine.NovaEngine.FlavorID)

	// Set scale type and related fields
	model.ScaleType, model.FixedScaleNodeCount, model.AutoScaleNodeCount, model.AutoScaleMinSize, model.AutoScaleMaxSize, diags = FlattenScaleSpec(nodeGroup.ScaleSpec)
	if diags.HasError() {
		return model, diags
	}

	// Set taints
	model.Taints, diags = FlattenTaints(ctx, nodeGroup.Taints)
	if diags.HasError() {
		return model, diags
	}

	// Set root disk info
	model.DiskType, model.DiskSize, diags = FlattenDiskType(nodeGroup.DiskType)
	if diags.HasError() {
		return model, diags
	}

	return model, diags
}
