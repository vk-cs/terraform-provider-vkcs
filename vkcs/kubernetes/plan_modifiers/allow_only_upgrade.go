package plan_modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func AllowOnlyUpgrade() planmodifier.String {
	return allowOnlyUpgradeModifier{}
}

type allowOnlyUpgradeModifier struct{}

func (m allowOnlyUpgradeModifier) Description(ctx context.Context) string {
	return m.MarkdownDescription(ctx)
}

func (m allowOnlyUpgradeModifier) MarkdownDescription(_ context.Context) string {
	return "Prevents downgrading to an older version. Only upgrades to newer versions are allowed."
}

func (m allowOnlyUpgradeModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Skip if state value is null (new resource creation)
	if req.StateValue.IsNull() {
		return
	}

	// Skip if plan value is unknown
	if req.PlanValue.IsUnknown() {
		return
	}

	// Skip if values are the same (no change)
	if req.StateValue.Equal(req.PlanValue) {
		return
	}

	oldVersion := req.StateValue.ValueString()
	newVersion := req.PlanValue.ValueString()

	// Simple string comparison for versions with same format "v1.2.3"
	// Because "v1.32.1" > "v1.31.4" lexicographically
	if newVersion < oldVersion {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Kubernetes version downgrade not allowed",
			fmt.Sprintf("Cannot downgrade from %s to %s. Only upgrades to newer versions are permitted.", oldVersion, newVersion),
		)
	}
}
