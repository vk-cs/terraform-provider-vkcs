package dataplatform

import (
	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func clusterStateRefreshFunc(client *gophercloud.ServiceClient, clusterID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, err := clusters.Get(client, clusterID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return c, string(clusterStatusDeleted), nil
			}
			return nil, "", err
		}

		return c, c.Status, nil
	}
}
