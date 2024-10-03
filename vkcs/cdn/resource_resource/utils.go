package resource_resource

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func OptionEnabled(option basetypes.ObjectValue) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if option.IsNull() || option.IsUnknown() {
		return false, diags
	}

	attributes := option.Attributes()
	enabledAttribute, ok := attributes["enabled"]
	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)
		return false, diags
	}

	enabledVal, ok := enabledAttribute.(basetypes.BoolValue)
	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`enabled expected to be basetypes.BoolValue, was: %T`, enabledAttribute))
		return false, diags
	}

	return enabledVal.ValueBool(), diags
}
