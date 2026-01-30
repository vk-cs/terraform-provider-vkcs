package schema_validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

var (
	_ validator.Map = (*KubernetesClusterLabelsValidator)(nil)
)

const (
	clusterLabelKeyMaxLen   = 63
	clusterLabelValueMaxLen = 63
)

type KubernetesClusterLabelsValidator struct{}

func (v KubernetesClusterLabelsValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v KubernetesClusterLabelsValidator) MarkdownDescription(ctx context.Context) string {
	return "Cluster labels maximum length of 63 characters for both keys and values"
}

func (v KubernetesClusterLabelsValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if util.IsNullOrUnknown(req.ConfigValue) {
		return
	}

	labelElements := req.ConfigValue.Elements()

	for key, value := range labelElements {
		// validate cluster label key length
		keyLength := len(key)
		if keyLength > clusterLabelKeyMaxLen {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Cluster label key too long",
				fmt.Sprintf("Key must be no more than %d characters, got %d characters", clusterLabelKeyMaxLen, keyLength),
			)
		}

		// validate cluster label value type
		strVal, ok := value.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Invalid cluster label value type",
				fmt.Sprintf("Expected string value, got %T", value),
			)
			continue
		}

		// validate cluster label value length
		valueLength := len(strVal.ValueString())
		if valueLength > clusterLabelValueMaxLen {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Cluster label value too long",
				fmt.Sprintf("Value must be no more than %d characters, got %d characters", clusterLabelValueMaxLen, valueLength),
			)
		}
	}
}
