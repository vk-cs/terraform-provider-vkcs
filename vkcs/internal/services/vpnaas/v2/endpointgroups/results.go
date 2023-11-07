package endpointgroups

import "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/endpointgroups"

func ExtractEndpointGroupInto(r endpointgroups.GetResult, v interface{}) error {
	return r.ExtractIntoStructPtr(v, "endpoint_group")
}
