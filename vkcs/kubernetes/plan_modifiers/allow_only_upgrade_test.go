package plan_modifiers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanUpgradeCluster(t *testing.T) {
	tests := []struct {
		name        string
		fromVersion string
		toVersion   string
		wantError   string
	}{
		// Allowed upgrades.
		{name: "PatchUpgradeWithMultiDigitPatch", fromVersion: "v1.31.4", toVersion: "v1.31.10"},
		{name: "PatchUpgrade", fromVersion: "v1.31.4", toVersion: "v1.31.6"},
		{name: "MinorUpgradeWithZeroPatch", fromVersion: "v1.31.4", toVersion: "v1.32.0"},
		{name: "MinorUpgrade", fromVersion: "v1.31.4", toVersion: "v1.32.1"},
		{name: "MultiDigitMinorUpgrade", fromVersion: "v1.9.4", toVersion: "v1.10.0"},

		// Not allowed.
		{name: "DowngradeToLowerPatch", fromVersion: "v1.31.10", toVersion: "v1.31.4", wantError: "only patch and minor upgrades are allowed"},
		{name: "DowngradeToLowerMinor", fromVersion: "v1.32.0", toVersion: "v1.31.4", wantError: "only patch and minor upgrades are allowed"},
		{name: "DowngradeThatLexicographicallyLooksBigger", fromVersion: "v1.31.4", toVersion: "v1.4.0", wantError: "only patch and minor upgrades are allowed"},
		{name: "SkipAMinorVersion", fromVersion: "v1.31.4", toVersion: "v1.33.0", wantError: "only patch and minor upgrades are allowed"},
		{name: "MajorUpgradeNotAllowed", fromVersion: "v1.31.4", toVersion: "v2.0.0", wantError: "only patch and minor upgrades are allowed"},

		// Invalid versions.
		{name: "InvalidFromVersion", fromVersion: "v1.31.4.5", toVersion: "v1.31.4", wantError: "current cluster version 'v1.31.4.5' is invalid"},
		{name: "InvalidToVersion", fromVersion: "v1.31.4", toVersion: "1.31.10", wantError: "target cluster version '1.31.10' is invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := canUpgradeCluster(tt.fromVersion, tt.toVersion)

			if tt.wantError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantError)
				return
			}

			assert.NoError(t, err)
		})
	}
}
