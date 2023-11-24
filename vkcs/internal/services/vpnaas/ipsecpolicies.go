package vpnaas

import "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"

func ExtractIPSecPolicyInto(r ipsecpolicies.GetResult, v interface{}) error {
	return r.ExtractIntoStructPtr(v, "ipsecpolicy")
}
