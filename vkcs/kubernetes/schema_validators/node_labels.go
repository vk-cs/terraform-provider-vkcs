package schema_validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/schema_validators/helpers"
)

var (
	_ validator.Map = (*KubernetesNodeLabelsValidator)(nil)
)

type KubernetesNodeLabelsValidator struct{}

func (v KubernetesNodeLabelsValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v KubernetesNodeLabelsValidator) MarkdownDescription(ctx context.Context) string {
	return "Node labels must follow Kubernetes label conventions"
}

func (v KubernetesNodeLabelsValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if util.IsNullOrUnknown(req.ConfigValue) {
		return
	}

	labelElements := req.ConfigValue.Elements()

	for key, value := range labelElements {
		// validate key using IsQualifiedName
		if err := helpers.IsQualifiedName(key); err != nil {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Invalid node group label key",
				err.Error(),
			)
		}

		// validate value
		strVal, ok := value.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Invalid node group label value type",
				fmt.Sprintf("Expected string value, got %T", value),
			)
			continue
		}

		// validate value using IsValidLabelValue
		if err := helpers.IsValidLabelValue(strVal.ValueString()); err != nil {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Invalid node group label value",
				err.Error(),
			)
		}
	}
}
