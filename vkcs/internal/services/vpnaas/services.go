package vpnaas

import "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/services"

func ExtractServiceInto(r services.GetResult, v interface{}) error {
	return r.ExtractIntoStructPtr(v, "vpnservice")
}
