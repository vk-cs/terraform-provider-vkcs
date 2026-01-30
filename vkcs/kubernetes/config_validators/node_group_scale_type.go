package config_validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	rkubengv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/resource_kubernetes_node_group_v2"
)

type ScaleTypeConfigValidator struct{}

func (v ScaleTypeConfigValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v ScaleTypeConfigValidator) MarkdownDescription(ctx context.Context) string {
	return "Validates that fixed_scale and auto_scale attributes are mutually exclusive based on scale_type"
}

func (v ScaleTypeConfigValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data rkubengv2.KubernetesNodeGroupV2Model

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scaleType := data.ScaleType.ValueString()

	switch scaleType {
	case "fixed_scale":
		if !util.IsNullOrUnknown(data.AutoScaleMinSize) {
			resp.Diagnostics.AddAttributeError(
				path.Root("auto_scale_min_size"),
				"Invalid Configuration",
				"auto_scale_min_size cannot be set when scale_type is 'fixed_scale'",
			)
		}
		if !util.IsNullOrUnknown(data.AutoScaleMaxSize) {
			resp.Diagnostics.AddAttributeError(
				path.Root("auto_scale_max_size"),
				"Invalid Configuration",
				"auto_scale_max_size cannot be set when scale_type is 'fixed_scale'",
			)
		}

		if util.IsNullOrUnknown(data.FixedScaleNodeCount) {
			resp.Diagnostics.AddAttributeError(
				path.Root("fixed_scale_node_count"),
				"Invalid Configuration",
				"fixed_scale_node_count is required when scale_type is 'fixed_scale'",
			)
		}

	case "auto_scale":
		if !util.IsNullOrUnknown(data.FixedScaleNodeCount) {
			resp.Diagnostics.AddAttributeError(
				path.Root("fixed_scale_node_count"),
				"Invalid Configuration",
				"fixed_scale_node_count cannot be set when scale_type is 'auto_scale'",
			)
		}

		if util.IsNullOrUnknown(data.AutoScaleMinSize) {
			resp.Diagnostics.AddAttributeError(
				path.Root("auto_scale_min_size"),
				"Invalid Configuration",
				"auto_scale_min_size is required when scale_type is 'auto_scale'",
			)
		}
		if util.IsNullOrUnknown(data.AutoScaleMaxSize) {
			resp.Diagnostics.AddAttributeError(
				path.Root("auto_scale_max_size"),
				"Invalid Configuration",
				"auto_scale_max_size is required when scale_type is 'auto_scale'",
			)
		}

		minSize := int(data.AutoScaleMinSize.ValueInt64())
		maxSize := int(data.AutoScaleMaxSize.ValueInt64())

		if minSize > maxSize {
			resp.Diagnostics.AddError(
				"Invalid Configuration",
				fmt.Sprintf("for auto_scale, condition 'auto_scale_min_size <= auto_scale_max_size' must be met (min=%d,  max=%d)", minSize, maxSize),
			)
		}

	default:
		resp.Diagnostics.AddAttributeError(
			path.Root("scale_type"),
			"Invalid Configuration",
			fmt.Sprintf("scale_type must be either 'fixed_scale' or 'auto_scale', got %s", scaleType),
		)

	}
}
