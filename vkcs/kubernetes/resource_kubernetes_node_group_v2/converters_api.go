package resource_kubernetes_node_group_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

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
			Size:    int(model.AutoScaleMinSize.ValueInt64()),
		}
	}

	// Build labels
	labels := make(map[string]string, len(model.Labels.Elements()))
	if !util.IsNullOrUnknown(model.Labels) {
		diags.Append(model.Labels.ElementsAs(ctx, &labels, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Build taints
	taints := make([]nodegroups.Taint, 0, len(model.Taints.Elements()))
	if !util.IsNullOrUnknown(model.Taints) {
		taintValues := make([]TaintsValue, 0, len(model.Taints.Elements()))
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
