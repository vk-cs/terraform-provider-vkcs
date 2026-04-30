package schema_validators

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

var (
	_ validator.String = (*KubernetesVersionValidator)(nil)
)

type KubernetesVersionValidator struct{}

func (v KubernetesVersionValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v KubernetesVersionValidator) MarkdownDescription(ctx context.Context) string {
	return "Kubernetes version must be in format vX.XX.XX where X is a number (e.g., v1.34.2)"
}

func (v KubernetesVersionValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if util.IsNullOrUnknown(req.ConfigValue) {
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

//nolint:staticcheck
func isValidClusterVersion(version string) error {
	// basic format validation: vX.XX.XX
	versionRegex := regexp.MustCompile(`^v(-?\d+)\.(-?\d+)\.(-?\d+)$`)
	matches := versionRegex.FindStringSubmatch(version)

	if matches == nil {
		return fmt.Errorf("Version '%s' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)", version)
	}

	major, err1 := strconv.Atoi(matches[1])
	minor, err2 := strconv.Atoi(matches[2])
	patch, err3 := strconv.Atoi(matches[3])

	if err1 != nil || err2 != nil || err3 != nil {
		return fmt.Errorf("Version components must be numbers, got '%s'", version)
	}

	if major < 1 {
		return fmt.Errorf("Major version of cluster version must be at least 1, got %d", major)
	}

	if minor < 0 {
		return fmt.Errorf("Minor version of cluster version cannot be negative, got %d", minor)
	}

	if patch < 0 {
		return fmt.Errorf("Patch version of cluster version cannot be negative, got %d", patch)
	}

	checkLeadingZeros := func(component string) bool {
		return component != "0" && strings.HasPrefix(component, "0")
	}

	if checkLeadingZeros(matches[1]) {
		return fmt.Errorf("Version '%s' contains leading zeros in version components; version components must not have leading zeros", version)
	}
	if checkLeadingZeros(matches[2]) {
		return fmt.Errorf("Version '%s' contains leading zeros in version components; version components must not have leading zeros", version)
	}
	if checkLeadingZeros(matches[3]) {
		return fmt.Errorf("Version '%s' contains leading zeros in version components; version components must not have leading zeros", version)
	}

	return nil
}
