package validators

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// CIDRValidator validates that a string is a valid CIDR notation
type CIDRValidator struct{}

func (v CIDRValidator) Description(ctx context.Context) string {
	return "string must be a valid CIDR notation (e.g., 192.168.1.0/24 or 2001:db8::/32)"
}

func (v CIDRValidator) MarkdownDescription(ctx context.Context) string {
	return "String must be a valid CIDR notation (e.g., `192.168.1.0/24` or `2001:db8::/32`)"
}

func (v CIDRValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if _, _, err := net.ParseCIDR(value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid CIDR notation",
			fmt.Sprintf("Value '%s' is not a valid CIDR notation: %v", value, err),
		)
	}
}
