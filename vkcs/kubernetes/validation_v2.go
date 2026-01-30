package kubernetes

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/valid"
)

// isUUID validates that a string is a valid UUID
func isUUID(val any, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected string for %q, got %T", key, val))
		return
	}

	if _, err := uuid.Parse(v); err != nil {
		errs = append(errs, fmt.Errorf("%q must be a valid UUID, got: %s", key, v))
	}

	return
}

// isClusterNameV2 validates cluster name according to DNS subdomain rules
func isClusterNameV2(val any, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected string for %q, got %T", key, val))
		return
	}

	if err := valid.ClusterNameV2(v); err != nil {
		errs = append(errs, err)
	}

	return
}

// validateInsecureRegistryURLV2 validates insecure registry URL format
func validateInsecureRegistryURLV2(i interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	val, ok := i.(string)
	if !ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid type",
			Detail:   fmt.Sprintf("Expected string, got %T", i),
		})
		return diags
	}

	// Basic URL validation - can be extended based on requirements
	if val == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Empty registry URL",
			Detail:   "Registry URL cannot be empty",
		})
	}

	// Simple regex to validate basic URL format
	// This can be made more sophisticated based on actual requirements
	urlRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-.]*(?::[0-9]+)?(/[a-zA-Z0-9_.-]*)*$`)
	if !urlRegex.MatchString(val) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Potentially invalid registry URL",
			Detail:   fmt.Sprintf("Registry URL '%s' may not be valid. Expected format: hostname[:port][/path]", val),
		})
	}

	return diags
}
