package validators

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"golang.org/x/exp/slices"
)

var _ validator.String = &dateFilterValidator{}

type dateFilterValidator struct {
	filters []string
}

func (v dateFilterValidator) Description(_ context.Context) string {
	return fmt.Sprintf("string value should be either RFC3339 formatted time or time filter in format `filter:time`, "+
		"where `filter` is one of [%s] and `time` is RFC3339 formatted time", strings.Join(v.filters, ", "))
}

func (v dateFilterValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v dateFilterValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	var parts []string
	var resultErr *multierror.Error

	if regexp.MustCompile("^" + strings.Join(v.filters, "|") + ":").Match([]byte(value)) {
		parts = strings.SplitN(value, ":", 2)
	} else {
		parts = []string{value}
	}

	if len(parts) == 2 {
		if !slices.Contains(v.filters, parts[0]) {
			resultErr = multierror.Append(resultErr, errors.New("invalid date filter"))
		}
		if _, err := time.Parse(time.RFC3339, parts[1]); err != nil {
			resultErr = multierror.Append(resultErr, err)
		}

	} else {
		if _, err := time.Parse(time.RFC3339, parts[0]); err != nil {
			resultErr = multierror.Append(resultErr, err)
		}
	}

	if resultErr.ErrorOrNil() != nil {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			req.Path,
			v.Description(ctx),
			value,
		))
	}
}

func DateFilter(filters ...string) validator.String {
	return dateFilterValidator{
		filters: filters,
	}
}
