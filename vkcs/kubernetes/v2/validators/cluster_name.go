package validators

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const (
	dnsSubdomainFmt      string = "^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"
	nameErrorMsg         string = "name must consist of lower case latin alphanumeric characters, '-', must start and end with latin alphanumeric character, and must not contain consecutive hyphens"
	clusterNameMaxLength int    = 25
	clusterNameMinLength int    = 3
)

var (
	dnsSubdomainRegexp = regexp.MustCompile("^" + dnsSubdomainFmt + "$")
	validCharsRegex    = regexp.MustCompile(`[^-a-z0-9]`)
	consecutiveHyphens = regexp.MustCompile(`--`) // check for two consecutive hyphens
)

// ClusterNameValidator validates cluster name according to DNS subdomain rules
type ClusterNameValidator struct{}

func (v ClusterNameValidator) Description(ctx context.Context) string {
	return "cluster name must consist of lower case latin alphanumeric characters, '-', must start and end with latin alphanumeric character, and must not contain consecutive hyphens"
}

func (v ClusterNameValidator) MarkdownDescription(ctx context.Context) string {
	return "Cluster name must consist of lower case latin alphanumeric characters, '-', must start and end with latin alphanumeric character, and must not contain consecutive hyphens"
}

func (v ClusterNameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	name := req.ConfigValue.ValueString()

	if err := isValidClusterName(name); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid cluster name",
			err.Error(),
		)
	}
}

func isValidClusterName(name string) error {
	if len(name) < clusterNameMinLength || len(name) > clusterNameMaxLength {
		return fmt.Errorf("name length must be no less than %d and no more than %d characters, got %d",
			clusterNameMinLength, clusterNameMaxLength, len(name))
	}

	// check for consecutive hyphens
	if consecutiveHyphens.MatchString(name) {
		return fmt.Errorf("consecutive hyphens (--) are not allowed")
	}

	// format validation using regex
	if !dnsSubdomainRegexp.MatchString(name) {
		// get invalid characters for more informative error message
		invalidChars := getInvalidChars(name)
		if len(invalidChars) > 0 {
			return fmt.Errorf("%s, got invalid characters: %v", nameErrorMsg, invalidChars)
		}

		return errors.New(nameErrorMsg)
	}

	return nil
}

// getInvalidChars returns unique invalid characters in the name
func getInvalidChars(name string) []string {
	return removeDuplicates(validCharsRegex.FindAllString(name, -1))
}

// removeDuplicates removes duplicates from string slice
func removeDuplicates(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}

	seen := make(map[string]struct{})
	result := make([]string, 0, len(slice))

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
