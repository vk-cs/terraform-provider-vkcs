package vkcs

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
)

// networkingSecgroupStateRefreshFuncDelete returns a special case resource.StateRefreshFunc to try to delete a secgroup.
func networkingSecgroupStateRefreshFuncDelete(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete vkcs_networking_secgroup %s", id)

		r, err := groups.Get(networkingClient, id).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted vkcs_networking_secgroup %s", id)
				return r, "DELETED", nil
			}

			return r, "ACTIVE", err
		}

		err = groups.Delete(networkingClient, id).ExtractErr()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted vkcs_networking_secgroup %s", id)
				return r, "DELETED", nil
			}
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				return r, "ACTIVE", nil
			}

			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] vkcs_networking_secgroup %s is still active", id)

		return r, "ACTIVE", nil
	}
}
