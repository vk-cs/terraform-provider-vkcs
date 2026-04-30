package schema_validators

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

var (
	_ validator.String = (*LbAllowedCIDRValidator)(nil)
)

type LbAllowedCIDRValidator struct{}

func (v LbAllowedCIDRValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v LbAllowedCIDRValidator) MarkdownDescription(ctx context.Context) string {
	return "String must be a valid CIDR notation (e.g., 192.168.1.0/24 or 2001:db8::/32)"
}

func (v LbAllowedCIDRValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if util.IsNullOrUnknown(req.ConfigValue) {
		return
	}

	cidr := req.ConfigValue.ValueString()
	if _, _, err := net.ParseCIDR(cidr); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid CIDR notation",
			fmt.Sprintf("Value '%s' is not a valid CIDR notation: %v", cidr, err),
		)
	}
}
