package networking

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	ifloatingips "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking/v2/floatingips"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type floatingIPExtended struct {
	floatingips.FloatingIP
	dns.FloatingIPDNSExt
	networking.SDNExt
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
		fip, err := ifloatingips.Get(client, fipID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return fip, "DELETED", nil
			}

			return nil, "", err
		}

		return fip, fip.Status, nil
	}
}
