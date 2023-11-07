package publicdns

import (
	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/publicdns/v2/zones"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func publicDNSZoneStateRefreshFunc(client *gophercloud.ServiceClient, zoneID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		zone, err := zones.Get(client, zoneID).Extract()

		if err != nil {
			if errutil.IsNotFound(err) {
				return zone, zoneStatusDeleted, nil
			}
			return nil, "", err
		}
		return zone, zone.Status, nil
	}
}
