package resource_kubernetes_node_group_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
)

// ToCreateOpts конвертирует KubernetesNodeGroupV2Model в nodegroups.CreateOpts
func ToCreateOpts(ctx context.Context, model *KubernetesNodeGroupV2Model, clusterID string) (*nodegroups.CreateOpts, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Build VM engine
	vmEngine := nodegroups.VMEngine{
		NovaEngine: nodegroups.NovaEngine{
			FlavorID: model.NodeFlavor.ValueString(),
		},
	}

	// Build scale specification
	scaleType := model.ScaleType.ValueString()
	scaleSpec := nodegroups.ScaleSpec{}

	switch scaleType {
	case "fixed_scale":
		scaleSpec.FixedScale = &nodegroups.FixedScale{
			Size: int(model.FixedScaleNodeCount.ValueInt64()),
		}
	case "auto_scale":
		scaleSpec.AutoScale = &nodegroups.AutoScale{
			MinSize: int(model.AutoScaleMinSize.ValueInt64()),
			MaxSize: int(model.AutoScaleMaxSize.ValueInt64()),
			Size:    int(model.AutoScaleNodeCount.ValueInt64()),
		}
	}

	// Build labels
	labels := make(map[string]string)
	if !model.Labels.IsNull() {
		elements := make(map[string]types.String, len(model.Labels.Elements()))
		diags.Append(model.Labels.ElementsAs(ctx, &elements, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for k, v := range elements {
			labels[k] = v.ValueString()
		}
	}

	// Build taints
	var taints []nodegroups.Taint
	if !model.Taints.IsNull() {
		var taintValues []TaintsValue
		diags.Append(model.Taints.ElementsAs(ctx, &taintValues, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, taint := range taintValues {
			taints = append(taints, nodegroups.Taint{
				Key:    taint.Key.ValueString(),
				Value:  taint.Value.ValueString(),
				Effect: taint.Effect.ValueString(),
			})
		}
	}

	diskTypeConfig := nodegroups.DiskType{
		CinderVolumeType: nodegroups.CinderVolumeType{
			Type: model.DiskType.ValueString(),
			Size: int(model.DiskSize.ValueInt64()),
		},
	}

	spec := nodegroups.NodeGroupSpec{
		Name:                 model.Name.ValueString(),
		VMEngine:             vmEngine,
		Zones:                []string{model.AvailabilityZone.ValueString()},
		ScaleSpec:            scaleSpec,
		Labels:               labels,
		Taints:               taints,
		ParallelUpgradeChunk: int(model.ParallelUpgradeChunk.ValueInt64()),
		DiskType:             diskTypeConfig,
	}

	// Build create options
	createOpts := nodegroups.CreateOpts{
		ClusterID: model.ClusterId.ValueString(),
		Spec:      spec,
	}

	return &createOpts, diags
}
