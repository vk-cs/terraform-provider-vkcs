package vkcs

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func databaseBackupStateRefreshFunc(client databaseClient, backupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		b, err := dbBackupGet(client, backupID).extract()
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
