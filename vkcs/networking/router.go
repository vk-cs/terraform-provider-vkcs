package networking

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func resourceNetworkingRouterStateRefreshFunc(client *gophercloud.ServiceClient, routerID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := routers.Get(client, routerID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return n, "DELETED", nil
			}

			return n, "", err
		}

		return n, n.Status, nil
	}
}
