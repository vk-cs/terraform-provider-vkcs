package db

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/backups"
)

func databaseBackupStateRefreshFunc(client *gophercloud.ServiceClient, backupID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		b, err := backups.Get(client, backupID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return b, "DELETED", nil
			}
			return nil, "", err
		}

		if b.Status == string(dbBackupStatusError) {
			return b, b.Status, fmt.Errorf("there was an error creating the database backup")
		}

		return b, b.Status, nil
	}
}
