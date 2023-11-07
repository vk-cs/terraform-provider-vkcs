package rules

import "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"

func ExtractSecurityGroupRuleInto(r rules.GetResult, v interface{}) error {
	return r.ExtractIntoStructPtr(v, "security_group_rule")
}
