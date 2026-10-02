package plan_modifiers

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver"
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
	return "Prevents downgrading or skipping cluster versions. Only patch and minor upgrades are allowed."
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

	if err := canUpgradeCluster(oldVersion, newVersion); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Kubernetes version upgrade",
			err.Error(),
		)
	}
}

func canUpgradeCluster(fromVersion, toVersion string) error {
	from, err := semver.NewVersion(fromVersion[1:])
	if err != nil {
		return fmt.Errorf("current cluster version '%s' is invalid", fromVersion)
	}

	to, err := semver.NewVersion(toVersion[1:])
	if err != nil {
		return fmt.Errorf("target cluster version '%s' is invalid", toVersion)
	}

	// Patch upgrade. Ex. from 1.31.4 => 1.31.6
	if from.Major() == to.Major() && from.Minor() == to.Minor() && from.Patch() < to.Patch() {
		return nil
	}

	// Minor upgrade. Ex. from 1.31.4 => 1.32.1
	if from.Major() == to.Major() && from.Minor()+1 == to.Minor() {
		return nil
	}

	return fmt.Errorf("cannot upgrade cluster version from '%s' to '%s': only patch and minor upgrades are allowed", fromVersion, toVersion)
}
