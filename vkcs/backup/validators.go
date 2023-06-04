package backup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = timeValidator{}

type timeValidator struct{}

func (v timeValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v timeValidator) MarkdownDescription(_ context.Context) string {
	return "value must be of valid time format hh:mm"
}

func (v timeValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue
	_, err := parseTime(value.ValueString())
	if err == nil {
		return
	}

	response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
		request.Path,
		v.Description(ctx),
		value.String(),
	))
}
