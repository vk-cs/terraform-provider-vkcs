package validators

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// KubernetesVersionValidator validates Kubernetes version format vX.XX.XX
type KubernetesVersionValidator struct{}

func (v KubernetesVersionValidator) Description(ctx context.Context) string {
	return "Kubernetes version must be in format vX.XX.XX where X is a number (e.g., v1.34.2)"
}

func (v KubernetesVersionValidator) MarkdownDescription(ctx context.Context) string {
	return "Kubernetes version must be in format `vX.XX.XX` where X is a number (e.g., `v1.34.2`)"
}

func (v KubernetesVersionValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	version := req.ConfigValue.ValueString()

	if err := isValidClusterVersion(version); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid cluster version",
			err.Error(),
		)
	}
}

func isValidClusterVersion(version string) error {
	// basic format validation: vX.XX.XX (теперь разрешаем знак минус)
	versionRegex := regexp.MustCompile(`^v(-?\d+)\.(-?\d+)\.(-?\d+)$`)
	matches := versionRegex.FindStringSubmatch(version)

	if matches == nil {
		return fmt.Errorf("version '%s' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)", version)
	}

	major, err1 := strconv.Atoi(matches[1])
	minor, err2 := strconv.Atoi(matches[2])
	patch, err3 := strconv.Atoi(matches[3])

	if err1 != nil || err2 != nil || err3 != nil {
		return fmt.Errorf("version components must be numbers, got '%s'", version)
	}

	if major < 1 {
		return fmt.Errorf("major version of cluster version must be at least 1, got %d", major)
	}

	if minor < 0 {
		return fmt.Errorf("minor version of cluster version cannot be negative, got %d", minor)
	}

	if patch < 0 {
		return fmt.Errorf("patch version of cluster version cannot be negative, got %d", patch)
	}

	checkLeadingZeros := func(component string) bool {
		return component != "0" && strings.HasPrefix(component, "0")
	}

	if checkLeadingZeros(matches[1]) {
		return fmt.Errorf("version '%s' contains leading zeros in version components; version components must not have leading zeros", version)
	}
	if checkLeadingZeros(matches[2]) {
		return fmt.Errorf("version '%s' contains leading zeros in version components; version components must not have leading zeros", version)
	}
	if checkLeadingZeros(matches[3]) {
		return fmt.Errorf("version '%s' contains leading zeros in version components; version components must not have leading zeros", version)
	}

	return nil
}
