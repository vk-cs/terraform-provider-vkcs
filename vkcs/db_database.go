package vkcs

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/databases"
)

func databaseDatabaseStateRefreshFunc(client *gophercloud.ServiceClient, dbmsID string, databaseName string, dbmsType string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pages, err := databases.List(client, dbmsID, dbmsType).AllPages()
		if err != nil {
			return nil, "", fmt.Errorf("unable to retrieve vkcs database databases: %s", err)
		}

		allDatabases, err := databases.ExtractDatabases(pages)
		if err != nil {
			return nil, "", fmt.Errorf("unable to extract vkcs database databases: %s", err)
		}

		for _, v := range allDatabases {
			if v.Name == databaseName {
				return v, "ACTIVE", nil
			}
		}

		return nil, "BUILD", nil
	}
}

func databaseDatabaseExists(client *gophercloud.ServiceClient, dbmsID string, databaseName string, dbmsType string) (bool, error) {
	var err error

	pages, err := databases.List(client, dbmsID, dbmsType).AllPages()
	if err != nil {
		return false, err
	}

	allDatabases, err := databases.ExtractDatabases(pages)
	if err != nil {
		return false, err
	}

	for _, v := range allDatabases {
		if v.Name == databaseName {
			return true, nil
		}
	}

	return false, err
}
