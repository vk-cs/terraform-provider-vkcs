package vkcs

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func resourceNetworkingRouterInterfaceStateRefreshFunc(networkingClient *gophercloud.ServiceClient, portID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := ports.Get(networkingClient, portID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return r, "DELETED", nil
			}

			return r, "", err
		}

		return r, r.Status, nil
	}
}

func resourceNetworkingRouterInterfaceDeleteRefreshFunc(networkingClient *gophercloud.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		routerID := d.Get("router_id").(string)
		routerInterfaceID := d.Id()

		log.Printf("[DEBUG] Attempting to delete vkcs_networking_router_interface %s", routerInterfaceID)

		removeOpts := routers.RemoveInterfaceOpts{
			SubnetID: d.Get("subnet_id").(string),
			PortID:   d.Get("port_id").(string),
		}

		if removeOpts.SubnetID != "" {
			// We need to make sure to only send subnet_id, because the port may have multiple
			// vkcs_networking_router_interface attached.
			removeOpts.PortID = ""
		}

		r, err := ports.Get(networkingClient, routerInterfaceID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted vkcs_networking_router_interface %s", routerInterfaceID)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		_, err = routers.RemoveInterface(networkingClient, routerID, removeOpts).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted vkcs_networking_router_interface %s", routerInterfaceID)
				return r, "DELETED", nil
			}
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				log.Printf("[DEBUG] vkcs_networking_router_interface %s is still in use", routerInterfaceID)
				return r, "ACTIVE", nil
			}

			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] vkcs_networking_router_interface %s is still active", routerInterfaceID)
		return r, "ACTIVE", nil
	}
}
