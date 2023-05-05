package vkcs

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	db "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/databases"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/users"
)

func extractDatabaseUserDatabases(v []interface{}) ([]db.CreateOpts, error) {
	Batch := make([]db.CreateOpts, len(v))
	for i, databaseName := range v {
		var C db.CreateOpts
		C.Name = databaseName.(string)
		Batch[i] = C
	}
	return Batch, nil
}

func flattenDatabaseUserDatabases(v []db.Database) []interface{} {
	databases := make([]interface{}, len(v))
	for i, db := range v {
		databases[i] = db.Name
	}
	return databases
}

func databaseUserStateRefreshFunc(client *gophercloud.ServiceClient, dbmsID string, userName string, dbmsType string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pages, err := users.List(client, dbmsID, dbmsType).AllPages()
		if err != nil {
			return nil, "", fmt.Errorf("unable to retrieve vkcs database users: %s", err)
		}

		allUsers, err := users.ExtractUsers(pages)
		if err != nil {
			return nil, "", fmt.Errorf("unable to extract vkcs database users: %s", err)
		}

		for _, v := range allUsers {
			if v.Name == userName {
				return v, "ACTIVE", nil
			}
		}

		return nil, "BUILD", nil
	}
}

func databaseUserExists(client *gophercloud.ServiceClient, dbmsID string, userName string, dbmsType string) (bool, *users.User, error) {
	var err error

	pages, err := users.List(client, dbmsID, dbmsType).AllPages()
	if err != nil {
		return false, nil, err
	}

	allUsers, err := users.ExtractUsers(pages)
	if err != nil {
		return false, nil, err
	}

	for _, v := range allUsers {
		if v.Name == userName {
			return true, &v, nil
		}
	}

	return false, nil, err
}
