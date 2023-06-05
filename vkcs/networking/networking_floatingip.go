package networking

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type floatingIPExtended struct {
	floatingips.FloatingIP
	dns.FloatingIPDNSExt
}

// networkingFloatingIPV2ID retrieves floating IP ID by the provided IP address.
func networkingFloatingIPV2ID(client *gophercloud.ServiceClient, floatingIP string) (string, error) {
	listOpts := floatingips.ListOpts{
		FloatingIP: floatingIP,
	}

	allPages, err := floatingips.List(client, listOpts).AllPages()
	if err != nil {
		return "", err
	}

	allFloatingIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return "", err
	}

	if len(allFloatingIPs) == 0 {
		return "", fmt.Errorf("there are no vkcs_networking_floatingip with %s IP", floatingIP)
	}
	if len(allFloatingIPs) > 1 {
		return "", fmt.Errorf("there are more than one vkcs_networking_floatingip with %s IP", floatingIP)
	}

	return allFloatingIPs[0].ID, nil
}

func networkingFloatingIPV2StateRefreshFunc(client *gophercloud.ServiceClient, fipID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		fip, err := floatingips.Get(client, fipID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return fip, "DELETED", nil
			}

			return nil, "", err
		}

		return fip, fip.Status, nil
	}
}
