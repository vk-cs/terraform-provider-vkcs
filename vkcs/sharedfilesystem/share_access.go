package sharedfilesystem

import (
	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/sharedfilesystem/v2/shares"
)

func sharedFilesystemShareAccessStateRefreshFunc(client *gophercloud.ServiceClient, shareID string, accessID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		access, err := shares.ListAccessRights(client, shareID).Extract()
		if err != nil {
			return nil, "", err
		}
		for _, v := range access {
			if v.ID == accessID {
				return v, v.State, nil
			}
		}
		return nil, "", gophercloud.ErrDefault404{}
	}
}
