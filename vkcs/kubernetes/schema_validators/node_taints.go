package schema_validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/schema_validators/helpers"
)

var (
	_ validator.Set = (*KubernetesNodeTaintsValidator)(nil)
)

type KubernetesNodeTaintsValidator struct{}

func (v KubernetesNodeTaintsValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v KubernetesNodeTaintsValidator) MarkdownDescription(ctx context.Context) string {
	return "Node taints must follow Kubernetes taints conventions"
}

//nolint:staticcheck
func (v KubernetesNodeTaintsValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if util.IsNullOrUnknown(req.ConfigValue) {
		return
	}

	taintElements := req.ConfigValue.Elements()

	for idx, elem := range taintElements {
		objValuable, ok := elem.(basetypes.ObjectValuable)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtSetValue(objValuable),
				"Invalid node taint format",
				fmt.Sprintf("Taint at index %d must be an object with 'key', 'value', and 'effect' string attributes. "+
					"Example: { key = \"key1\", value = \"value1\", effect = \"NoSchedule\" }", idx),
			)
			continue
		}

		taintPath := req.Path.AtSetValue(objValuable)

		objValue, objDiags := objValuable.ToObjectValue(ctx)
		if objDiags.HasError() {
			resp.Diagnostics.Append(objDiags...)
			continue
		}

		attrs := objValue.Attributes()

		// Validate key
		keyAttr, keyExists := attrs["key"]
		if !keyExists {
			resp.Diagnostics.AddAttributeError(
				taintPath.AtName("key"),
				"Missing taint key",
				fmt.Sprintf("Taint at index %d must have a 'key' attribute", idx),
			)
			continue
		}

		keyStr, keyOk := keyAttr.(basetypes.StringValue)
		if !keyOk {
			resp.Diagnostics.AddAttributeError(
				taintPath.AtName("key"),
				"Invalid taint key type",
				fmt.Sprintf("Taint at index %d: expected string value for 'key', got %T", idx, keyAttr),
			)
			continue
		}

		key := keyStr.ValueString()
		if err := v.isValidTaintKey(key); err != nil {
			resp.Diagnostics.AddAttributeError(
				taintPath.AtName("key"),
				"Invalid taint key",
				fmt.Sprintf("Taint at index %d with key '%s': %s", idx, key, err.Error()),
			)
		}

		// Validate value
		valueAttr, valueExists := attrs["value"]
		if !valueExists {
			resp.Diagnostics.AddAttributeError(
				taintPath.AtName("value"),
				"Missing taint value",
				fmt.Sprintf("Taint at index %d must have a 'value' attribute", idx),
			)
			continue
		}

		valueStr, valueOk := valueAttr.(basetypes.StringValue)
		if !valueOk {
			resp.Diagnostics.AddAttributeError(
				taintPath.AtName("value"),
				"Invalid taint value type",
				fmt.Sprintf("Taint at index %d: expected string value for 'value', got %T", idx, valueAttr),
			)
			continue
		}

		value := valueStr.ValueString()
		if err := v.isValidTaintValue(value); err != nil {
			resp.Diagnostics.AddAttributeError(
				taintPath.AtName("value"),
				"Invalid taint value",
				fmt.Sprintf("Taint at index %d with value '%s': %s", idx, value, err.Error()),
			)
		}

		effectAttr, effectExists := attrs["effect"]
		if !effectExists {
			resp.Diagnostics.AddAttributeError(
				taintPath.AtName("effect"),
				"Missing taint effect",
				fmt.Sprintf("Taint at index %d must have an 'effect' attribute", idx),
			)
			continue
		}

		effectStr, effectOk := effectAttr.(basetypes.StringValue)
		if !effectOk {
			resp.Diagnostics.AddAttributeError(
				taintPath.AtName("effect"),
				"Invalid taint effect type",
				fmt.Sprintf("Taint at index %d: expected string value for 'effect', got %T", idx, effectAttr),
			)
			continue
		}

		effect := effectStr.ValueString()
		if err := v.isValidTaintEffect(effect); err != nil {
			resp.Diagnostics.AddAttributeError(
				taintPath.AtName("effect"),
				"Invalid taint effect",
				fmt.Sprintf("Taint at index %d with effect '%s': %s", idx, effect, err.Error()),
			)
		}
	}
}

func (v KubernetesNodeTaintsValidator) isValidTaintKey(key string) error {
	return helpers.IsQualifiedName(key)
}

func (v KubernetesNodeTaintsValidator) isValidTaintValue(value string) error {
	return helpers.IsValidLabelValue(value)
}

//nolint:staticcheck
func (v KubernetesNodeTaintsValidator) isValidTaintEffect(effect string) error {
	// Valid taint effects in Kubernetes
	validEffects := []string{
		"NoSchedule",
		"PreferNoSchedule",
		"NoExecute",
	}

	for _, validEffect := range validEffects {
		if effect == validEffect {
			return nil
		}
	}

	return fmt.Errorf("Taint effect must be one of %s, got %s", strings.Join(validEffects, ", "), effect)
}
