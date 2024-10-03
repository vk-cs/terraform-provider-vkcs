package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.ConfigValidator = (*conflictingEnabledValidator)(nil)

func ConflictingEnabled(expressions ...path.Expression) resource.ConfigValidator {
	return conflictingEnabledValidator{
		pathExpressions: expressions,
	}
}

// conflictingEnabledValidator is the underlying struct implementing ConflictsWith.
type conflictingEnabledValidator struct {
	pathExpressions path.Expressions
}

func (v conflictingEnabledValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v conflictingEnabledValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("These attributes cannot be enabled simultaneously: %s", v.pathExpressions)
}

func (v conflictingEnabledValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v conflictingEnabledValidator) Validate(ctx context.Context, config tfsdk.Config) diag.Diagnostics {
	var configuredPaths path.Paths
	var diags diag.Diagnostics

	for _, expression := range v.pathExpressions {
		matchedPaths, d := config.PathMatches(ctx, expression)
		diags.Append(d...)
		if d.HasError() {
			continue
		}

		for _, matchedPath := range matchedPaths {
			var value attr.Value
			d = config.GetAttribute(ctx, matchedPath, &value)
			diags.Append(d...)
			if d.HasError() {
				continue
			}

			if value.IsNull() || value.IsUnknown() {
				continue
			}

			if v, ok := value.(basetypes.BoolValue); ok {
				if v.ValueBool() {
					configuredPaths.Append(matchedPath)
				}
				continue
			}

			v, ok := value.(basetypes.ObjectValuable)
			if !ok {
				diags.AddAttributeError(
					matchedPath,
					"Attribute Wrong Type",
					fmt.Sprintf(`value expected to be either basetypes.BoolValue or basetypes.ObjectValuable, was: %T`, value))
				continue
			}

			obj, d := v.ToObjectValue(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			enabledAttribute, ok := obj.Attributes()["enabled"]
			if !ok {
				diags.AddError(
					"Attribute Missing",
					`enabled is missing from object`)
				continue
			}

			enabledVal, ok := enabledAttribute.(basetypes.BoolValue)
			if !ok {
				diags.AddError(
					"Attribute Wrong Type",
					fmt.Sprintf(`enabled expected to be basetypes.BoolValue, was: %T`, enabledAttribute))
				continue
			}

			if enabledVal.ValueBool() {
				configuredPaths.Append(matchedPath)
			}
		}
	}

	if len(configuredPaths) > 1 {
		diags.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
			configuredPaths[0],
			v.Description(ctx),
		))
	}

	return diags
}
