package vkcs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/external"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/provider"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/pagination"
)

type PrivateDNSDomainExt struct {
	PrivateDNSDomain string `json:"private_dns_domain,omitempty"`
}

type ServicesAccessExt struct {
	ServicesAccess *bool `json:"enable_shadow_port,omitempty"`
}

type networkExtended struct {
	networks.Network
	external.NetworkExternalExt
	portsecurity.PortSecurityExt
	dns.NetworkDNSExt
	policies.QoSPolicyExt
	provider.NetworkProviderExt
	PrivateDNSDomainExt
	ServicesAccessExt
}

// networkingNetworkID retrieves network ID by the provided name.
func networkingNetworkID(d *schema.ResourceData, meta interface{}, networkName string) (string, error) {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return "", fmt.Errorf("error creating VKCS network client: %s", err)
	}

	opts := networks.ListOpts{Name: networkName}
	pager := networks.List(networkingClient, opts)
	networkID := ""

	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, err := networks.ExtractNetworks(page)
		if err != nil {
			return false, err
		}

		for _, n := range networkList {
			if n.Name == networkName {
				networkID = n.ID
				return false, nil
			}
		}

		return true, nil
	})

	return networkID, err
}

// networkingNetworkName retrieves network name by the provided ID.
func networkingNetworkName(d *schema.ResourceData, meta interface{}, networkID string) (string, error) {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return "", fmt.Errorf("error creating VKCS network client: %s", err)
	}

	opts := networks.ListOpts{ID: networkID}
	pager := networks.List(networkingClient, opts)
	networkName := ""

	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, err := networks.ExtractNetworks(page)
		if err != nil {
			return false, err
		}

		for _, n := range networkList {
			if n.ID == networkID {
				networkName = n.Name
				return false, nil
			}
		}

		return true, nil
	})

	return networkName, err
}

func resourceNetworkingNetworkStateRefreshFunc(client *gophercloud.ServiceClient, networkID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := networks.Get(client, networkID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return n, "DELETED", nil
			}
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				return n, "ACTIVE", nil
			}

			return n, "", err
		}

		return n, n.Status, nil
	}
}
