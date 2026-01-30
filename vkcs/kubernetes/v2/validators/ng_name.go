package validators

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	ngNameBlackList     = []string{"master"}
	ngNameMaxLength int = 25
	ngNameMinLength int = 3
)

// NodeGroupNameValidator validates node group name according to DNS subdomain rules
type NodeGroupNameValidator struct{}

func (v NodeGroupNameValidator) Description(ctx context.Context) string {
	return "node group name must consist of lower case latin alphanumeric characters, '-', must start and end with latin alphanumeric character, must not contain consecutive hyphens, and cannot be 'master'"
}

func (v NodeGroupNameValidator) MarkdownDescription(ctx context.Context) string {
	return "Node group name must consist of lower case latin alphanumeric characters, '-', must start and end with latin alphanumeric character, must not contain consecutive hyphens, and cannot be 'master'"
}

func (v NodeGroupNameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if err := isValidNodeGroupName(value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid node group name",
			err.Error(),
		)
	}
}

func isValidNodeGroupName(ngName string) error {
	ngName = strings.TrimSpace(ngName)

	if slices.Contains(ngNameBlackList, ngName) {
		return errors.New("name can't be a word 'master'")
	}

	if len(ngName) < ngNameMinLength || len(ngName) > ngNameMaxLength {
		return fmt.Errorf("name length must be no less than %d and no more than %d characters, got %d", ngNameMinLength, ngNameMaxLength, len(ngName))
	}

	if consecutiveHyphens.MatchString(ngName) {
		return fmt.Errorf("consecutive hyphens (--) are not allowed")
	}

	if !dnsSubdomainRegexp.MatchString(ngName) {
		invalidChars := getInvalidChars(ngName)
		if len(invalidChars) > 0 {
			return fmt.Errorf("%s, got invalid characters: %v", nameErrorMsg, invalidChars)
		}

		return errors.New(nameErrorMsg)
	}

	return nil
}
