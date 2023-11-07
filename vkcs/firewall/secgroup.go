package firewall

import (
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	igroups "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/firewall/v2/groups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type securityGroupExtended struct {
	groups.SecGroup
	networking.SDNExt
}

// networkingSecgroupStateRefreshFuncDelete returns a special case retry.StateRefreshFunc to try to delete a secgroup.
func networkingSecgroupStateRefreshFuncDelete(networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete vkcs_networking_secgroup %s", id)

		r, err := igroups.Get(networkingClient, id).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				log.Printf("[DEBUG] Successfully deleted vkcs_networking_secgroup %s", id)
				return r, "DELETED", nil
			}

			return r, "ACTIVE", err
		}

		err = igroups.Delete(networkingClient, id).ExtractErr()
		if err != nil {
			if errutil.IsNotFound(err) {
				log.Printf("[DEBUG] Successfully deleted vkcs_networking_secgroup %s", id)
				return r, "DELETED", nil
			}

			if errutil.Is(err, 409) {
				return r, "ACTIVE", nil
			}

			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] vkcs_networking_secgroup %s is still active", id)

		return r, "ACTIVE", nil
	}
}
