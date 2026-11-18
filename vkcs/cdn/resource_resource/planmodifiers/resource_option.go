package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ planmodifier.Object = resourceOptionPlanModifier{}

func ResourceOption() planmodifier.Object {
	return resourceOptionPlanModifier{}
}

// ResourceOptionWithState preserves prior state when an omitted optional+computed
// object is marked unknown by the framework. This is needed for nested list
// attributes, which can otherwise cause a non-empty plan after apply (framework
// issue #898). The resource option behavior is then applied to known plans.
func ResourceOptionWithState() planmodifier.Object {
	return resourceOptionWithStatePlanModifier{}
}

type resourceOptionWithStatePlanModifier struct{}

func (m resourceOptionWithStatePlanModifier) Description(ctx context.Context) string {
	return "Preserves the prior state for an omitted CDN resource option."
}

func (m resourceOptionWithStatePlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m resourceOptionWithStatePlanModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.ConfigValue.IsNull() && req.PlanValue.IsUnknown() {
		resp.PlanValue = req.StateValue
		return
	}

	resourceOptionPlanModifier{}.PlanModifyObject(ctx, req, resp)
}

type resourceOptionPlanModifier struct{}

func (m resourceOptionPlanModifier) Description(_ context.Context) string {
	return ""
}

func (m resourceOptionPlanModifier) MarkdownDescription(_ context.Context) string {
	return ""
}

func (m resourceOptionPlanModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.ConfigValue.IsUnknown() {
		return
	}

	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	planAttrs := req.PlanValue.Attributes()
	enabledAttr, ok := planAttrs["enabled"]
	if !ok {
		resp.Diagnostics.AddError(
			"Attribute Missing",
			`enabled is missing from object`)
		return
	}

	valueSet := false
	for name, value := range planAttrs {
		if name != "enabled" && !value.IsUnknown() && !value.IsNull() {
			valueSet = true
			break
		}
	}

	if enabledAttr.IsNull() || enabledAttr.IsUnknown() && valueSet {
		enabledAttr = types.BoolValue(true)
		planAttrs["enabled"] = enabledAttr
	}

	enabled, ok := enabledAttr.(basetypes.BoolValue)
	if !ok {
		resp.Diagnostics.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`enabled expected to be basetypes.BoolValue, was: %T`, enabledAttr))
	}

	if !enabled.ValueBool() {
		stateAttrs := req.StateValue.Attributes()
		configAttrs := req.ConfigValue.Attributes()

		for name, value := range planAttrs {
			if name == "enabled" || !value.IsUnknown() {
				continue
			}

			if configAttr, ok := configAttrs[name]; ok && configAttr.IsUnknown() {
				continue
			}

			if stateAttr, ok := stateAttrs[name]; ok && !stateAttr.IsNull() {
				planAttrs[name] = stateAttr
			}
		}
	}

	newPlanValue, diags := types.ObjectValue(req.PlanValue.AttributeTypes(ctx), planAttrs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = newPlanValue
}
