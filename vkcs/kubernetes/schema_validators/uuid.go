package schema_validators

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ validator.String = (*UUIDValidator)(nil)
)

type UUIDValidator struct{}

func (v UUIDValidator) Description(ctx context.Context) string {
	return "string must be a valid UUID"
}

func (v UUIDValidator) MarkdownDescription(ctx context.Context) string {
	return "String must be a valid UUID"
}

func (v UUIDValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if _, err := uuid.Parse(value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid UUID",
			fmt.Sprintf("String '%s' is not a valid UUID", value),
		)
	}
}
