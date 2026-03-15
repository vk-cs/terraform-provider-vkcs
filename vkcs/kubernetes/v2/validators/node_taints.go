package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/v2/validators/helpers"
)

// KubernetesTaintKeyValidator validates Kubernetes taint key
type KubernetesTaintKeyValidator struct{}

func (v KubernetesTaintKeyValidator) Description(ctx context.Context) string {
	return "taint key must be non-empty, max 253 chars, consist of alphanumeric characters, '-', '_', or '.', and start/end with alphanumeric character"
}

func (v KubernetesTaintKeyValidator) MarkdownDescription(ctx context.Context) string {
	return "Taint key must be non-empty, max 253 chars, consist of alphanumeric characters, '-', '_', or '.', and start/end with alphanumeric character"
}

func (v KubernetesTaintKeyValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if err := v.isValidTaintKey(value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid taint key",
			err.Error(),
		)
	}
}

func (v KubernetesTaintKeyValidator) isValidTaintKey(key string) error {
	return helpers.IsQualifiedName(key)
}

// KubernetesTaintValueValidator validates Kubernetes taint value
type KubernetesTaintValueValidator struct{}

func (v KubernetesTaintValueValidator) Description(ctx context.Context) string {
	return "taint value can be empty or max 63 chars, consist of alphanumeric characters, '-', '_', or '.', and start/end with alphanumeric character if not empty"
}

func (v KubernetesTaintValueValidator) MarkdownDescription(ctx context.Context) string {
	return "Taint value can be empty or max 63 chars, consist of alphanumeric characters, '-', '_', or '.', and start/end with alphanumeric character if not empty"
}

func (v KubernetesTaintValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if err := v.isValidTaintValue(value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid taint value",
			err.Error(),
		)
	}
}

func (v KubernetesTaintValueValidator) isValidTaintValue(value string) error {
	return helpers.IsValidLabelValue(value)
}
