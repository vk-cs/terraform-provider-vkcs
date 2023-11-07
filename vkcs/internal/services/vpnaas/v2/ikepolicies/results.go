package ikepolicies

import "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ikepolicies"

func ExtractIKEPolicyInto(r ikepolicies.GetResult, v interface{}) error {
	return r.ExtractIntoStructPtr(v, "ikepolicy")
}
