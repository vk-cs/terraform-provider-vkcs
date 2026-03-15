package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/segmentio/ksuid"
)

// KSUIDValidator проверяет, что строка является валидным KSUID
type KSUIDValidator struct{}

func (v KSUIDValidator) Description(ctx context.Context) string {
	return "string must be a valid KSUID"
}

func (v KSUIDValidator) MarkdownDescription(ctx context.Context) string {
	return "String must be a valid KSUID"
}

func (v KSUIDValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
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
