package schema_validators

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

var (
	_ validator.String = (*LbAllowedCIDRValidator)(nil)
	_ validator.Set    = (*LbAllowedCIDRsUniquenessValidator)(nil)
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

// LbAllowedCIDRsUniquenessValidator checks that loadbalancer_allowed_cidrs does
// not contain several entries that belong to the same subnet. Since the backend
// stores CIDRs in their network form (host bits are ignored), entries such as
// 192.168.10.1/24 and 192.168.10.2/24 would both be stored as 192.168.10.0/24
// and are therefore rejected as duplicates.
type LbAllowedCIDRsUniquenessValidator struct{}

func (v LbAllowedCIDRsUniquenessValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v LbAllowedCIDRsUniquenessValidator) MarkdownDescription(ctx context.Context) string {
	return "Each CIDR must belong to a distinct subnet; entries from the same subnet are not allowed"
}

//nolint:staticcheck
func (v LbAllowedCIDRsUniquenessValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if util.IsNullOrUnknown(req.ConfigValue) {
		return
	}

	var cidrs []string
	resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &cidrs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, conflict := range findConflictingCIDRs(cidrs) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"CIDRs from the same subnet",
			conflict,
		)
	}
}

// findConflictingCIDRs groups the given CIDRs by their normalized network
// address and returns a message for every group that contains more than one
// entry. Invalid CIDR strings are skipped: their format is validated separately
// by LbAllowedCIDRValidator.
func findConflictingCIDRs(cidrs []string) []string {
	groups := make(map[string][]string)
	for _, cidr := range cidrs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}

		network := ipNet.String()
		groups[network] = append(groups[network], cidr)
	}

	var conflicts []string
	for network, group := range groups {
		if len(group) < 2 {
			continue
		}

		sort.Strings(group)
		conflicts = append(conflicts, fmt.Sprintf(
			"CIDRs %s belong to the same subnet and would all be stored as %q",
			quoteStrings(group), network,
		))
	}

	sort.Strings(conflicts)
	return conflicts
}

func quoteStrings(items []string) string {
	quoted := make([]string, len(items))
	for i, s := range items {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	return strings.Join(quoted, ", ")
}
