package backup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ validator.Object = objectLockConfigValidator{}

type objectLockConfigValidator struct{}

func (objectLockConfigValidator) Description(context.Context) string {
	return "retention fields must match the configured Object Lock mode"
}

func (v objectLockConfigValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (objectLockConfigValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var config PlanResourceObjectLockModel
	resp.Diagnostics.Append(req.ConfigValue.As(ctx, &config, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateObjectLockSettings(req.Path, config)...)
}

func validateObjectLockSettings(attributePath path.Path, config PlanResourceObjectLockModel) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	if config.Type.IsNull() || config.Type.IsUnknown() {
		return diagnostics
	}

	daysConfigured := !config.Days.IsNull()
	retainUntilConfigured := !config.RetainUntil.IsNull()
	switch config.Type.ValueString() {
	case "compliance", "governance":
		if daysConfigured == retainUntilConfigured {
			diagnostics.AddAttributeError(
				attributePath,
				"Invalid Object Lock retention configuration",
				"For compliance or governance mode, set exactly one of `days` or `retain_until`.",
			)
		}
	case "legalhold":
		if daysConfigured {
			diagnostics.AddAttributeError(
				attributePath.AtName("days"),
				"Invalid Object Lock retention configuration",
				"`days` must not be set when Object Lock type is `legalhold`.",
			)
		}
		if retainUntilConfigured {
			diagnostics.AddAttributeError(
				attributePath.AtName("retain_until"),
				"Invalid Object Lock retention configuration",
				"`retain_until` must not be set when Object Lock type is `legalhold`.",
			)
		}
	}

	return diagnostics
}
