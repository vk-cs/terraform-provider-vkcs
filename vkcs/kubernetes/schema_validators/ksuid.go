package schema_validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/segmentio/ksuid"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

var (
	_ validator.String = (*KSUIDValidator)(nil)
)

type KSUIDValidator struct{}

func (v KSUIDValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v KSUIDValidator) MarkdownDescription(ctx context.Context) string {
	return "String must be a valid KSUID"
}

func (v KSUIDValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if util.IsNullOrUnknown(req.ConfigValue) {
		return
	}

	value := req.ConfigValue.ValueString()

	if _, err := ksuid.Parse(value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid KSUID",
			fmt.Sprintf("String '%s' is not a valid KSUID", value),
		)
	}
}
