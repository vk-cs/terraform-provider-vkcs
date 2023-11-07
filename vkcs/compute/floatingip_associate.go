package compute

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	nfloatingips "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	iservers "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/servers"
)

func ParseComputeFloatingIPAssociateID(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("unable to determine floating ip association ID")
	}

	floatingIP := idParts[0]
	instanceID := idParts[1]
	fixedIP := idParts[2]

	return floatingIP, instanceID, fixedIP, nil
}

func computeFloatingIPAssociateNetworkExists(networkClient *gophercloud.ServiceClient, floatingIP string) (bool, error) {
	listOpts := nfloatingips.ListOpts{
		FloatingIP: floatingIP,
	}
	allPages, err := nfloatingips.List(networkClient, listOpts).AllPages()
	if err != nil {
		return false, err
	}

	allFips, err := nfloatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return false, err
	}

	if len(allFips) > 1 {
		return false, fmt.Errorf("more than one floating IP with %s address found", floatingIP)
	}

	if len(allFips) == 0 {
		return false, nil
	}

	return true, nil
}

func computeFloatingIPAssociateComputeExists(computeClient *gophercloud.ServiceClient, floatingIP string) (bool, error) {
	// If the Network API isn't available, fall back to the deprecated Compute API.
	allPages, err := floatingips.List(computeClient).AllPages()
	if err != nil {
		return false, err
	}

	allFips, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return false, err
	}

	for _, f := range allFips {
		if f.IP == floatingIP {
			return true, nil
		}
	}

	return false, nil
}

func computeFloatingIPAssociateCheckAssociation(
	computeClient *gophercloud.ServiceClient, instanceID, floatingIP string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := iservers.Get(computeClient, instanceID).Extract()
		if err != nil {
			return instance, "", err
		}

		var associated bool
		for _, networkAddresses := range instance.Addresses {
			for _, element := range networkAddresses.([]interface{}) {
				address := element.(map[string]interface{})
				if address["OS-EXT-IPS:type"] == "floating" && address["addr"] == floatingIP {
					associated = true
				}
			}
		}

		if associated {
			return instance, "ASSOCIATED", nil
		}

		return instance, "NOT_ASSOCIATED", nil
	}
}
