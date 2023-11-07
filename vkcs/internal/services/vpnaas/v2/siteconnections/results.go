package siteconnections

import "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/siteconnections"

func ExtractConnectionInto(r siteconnections.GetResult, v interface{}) error {
	return r.ExtractIntoStructPtr(v, "ipsec_site_connection")
}
