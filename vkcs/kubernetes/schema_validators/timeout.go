package schema_validators

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

var (
	_ validator.String = (*TimeoutValidator)(nil)
)

type TimeoutValidator struct{}

func (v TimeoutValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v TimeoutValidator) MarkdownDescription(ctx context.Context) string {
	return "Timeout must be a valid duration (e.g., '60m', '2h', '1h30m', '30s')"
}

func (v TimeoutValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if util.IsNullOrUnknown(req.ConfigValue) {
		return
	}

	value := req.ConfigValue.ValueString()

	if value == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid timeout",
			"Timeout cannot be empty",
		)
		return
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid timeout format",
			fmt.Sprintf("'%s' is not a valid duration format for timeout. Expected format like '60m', '2h', '1h30m', '30s'", value),
		)
		return
	}

	if duration <= 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid timeout value",
			"Timeout must be positive",
		)
		return
	}
}
