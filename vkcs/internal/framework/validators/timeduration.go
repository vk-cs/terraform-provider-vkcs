package validators

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = timeDurationValidator{}

type timeDurationValidator struct{}

func (v timeDurationValidator) Description(_ context.Context) string {
	return `should be a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`
}

func (v timeDurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v timeDurationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if _, err := time.ParseDuration(value); err != nil {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(req.Path, v.Description(ctx), value))
		return
	}
}

func TimeDuration() validator.String {
	return timeDurationValidator{}
}
