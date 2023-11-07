package networking

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	irouters "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking/v2/routers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type routerExtended struct {
	routers.Router
	networking.SDNExt
}

func resourceNetworkingRouterStateRefreshFunc(client *gophercloud.ServiceClient, routerID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := irouters.Get(client, routerID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return n, "DELETED", nil
			}

			return n, "", err
		}

		return n, n.Status, nil
	}
}
