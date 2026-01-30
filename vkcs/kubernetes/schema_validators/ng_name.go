package schema_validators

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

var (
	_ validator.String = (*NodeGroupNameValidator)(nil)
)

var (
	ngNameBlackList     = []string{"master"}
	ngNameMaxLength int = 25
	ngNameMinLength int = 3
)

type NodeGroupNameValidator struct{}

func (v NodeGroupNameValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v NodeGroupNameValidator) MarkdownDescription(ctx context.Context) string {
	return "Node group name must consist of lower case Latin alphanumeric characters, '-', must start and end with Latin alphanumeric character, must not contain consecutive hyphens, and cannot be 'master'"
}

func (v NodeGroupNameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if util.IsNullOrUnknown(req.ConfigValue) {
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

//nolint:staticcheck
func isValidNodeGroupName(ngName string) error {
	ngName = strings.TrimSpace(ngName)

	if slices.Contains(ngNameBlackList, ngName) {
		return errors.New("Node group name can't be a word 'master'")
	}

	if len(ngName) < ngNameMinLength || len(ngName) > ngNameMaxLength {
		return fmt.Errorf("Node group name length must be no less than %d and no more than %d characters, got %d", ngNameMinLength, ngNameMaxLength, len(ngName))
	}

	if consecutiveHyphens.MatchString(ngName) {
		return fmt.Errorf("Consecutive hyphens (--) are not allowed")
	}

	if !dnsSubdomainRegexp.MatchString(ngName) {
		invalidChars := getInvalidChars(ngName)
		if len(invalidChars) > 0 {
			return fmt.Errorf("Node group %s, got invalid characters: %v", nameErrorMsg, invalidChars)
		}

		return errors.New("Node group " + nameErrorMsg)
	}

	return nil
}
