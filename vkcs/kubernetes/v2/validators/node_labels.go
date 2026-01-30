package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/v2/validators/helpers"
)

// KubernetesLabelsValidator validates Kubernetes labels map
type KubernetesLabelsValidator struct{}

func (v KubernetesLabelsValidator) Description(ctx context.Context) string {
	return "labels must follow kubernetes label conventions: keys have optional prefix and name separated by '/', values must be valid label values"
}

func (v KubernetesLabelsValidator) MarkdownDescription(ctx context.Context) string {
	return "labels must follow kubernetes label conventions: keys have optional prefix and name separated by '/', values must be valid label values"
}

func (v KubernetesLabelsValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()

	for key, value := range elements {
		// validate key using IsQualifiedName
		if err := helpers.IsQualifiedName(key); err != nil {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Invalid label key",
				err.Error(),
			)
		}

		// validate value
		var valueStr string
		if strVal, ok := value.(types.String); ok {
			valueStr = strVal.ValueString()
		} else {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Invalid label value type",
				fmt.Sprintf("expected string value, got %T", value),
			)
			continue
		}

		// validate value using IsValidLabelValue
		if err := helpers.IsValidLabelValue(valueStr); err != nil {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Invalid label value",
				err.Error(),
			)
		}
	}
}
