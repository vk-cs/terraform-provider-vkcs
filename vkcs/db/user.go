package db

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	db "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/databases"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/users"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func extractDatabaseUserDatabases(v []string) ([]db.CreateOpts, error) {
	Batch := make([]db.CreateOpts, len(v))
	for i, databaseName := range v {
		var C db.CreateOpts
		C.Name = databaseName
		Batch[i] = C
	}
	return Batch, nil
}

func flattenDatabaseUserDatabases(v []db.Database) []string {
	databases := make([]string, len(v))
	for i, db := range v {
		databases[i] = db.Name
	}
	return databases
}

func databaseUserStateRefreshFunc(client *gophercloud.ServiceClient, dbmsID string, userName string, dbmsType string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result := users.Get(client, dbmsID, userName, dbmsType)
		if result.Err != nil {
			if errutil.IsNotFound(result.Err) {
				return nil, "BUILD", nil
			}
			return nil, "", fmt.Errorf("unable to retrieve vkcs database user: %s", result.Err)
		}

		userObj, err := result.Extract()
		if err != nil {
			return nil, "", fmt.Errorf("unable to extract vkcs database user: %s", err)
		}
		if userObj == nil {
			return nil, "BUILD", nil
		}

		return userObj, "ACTIVE", nil
	}
}

func databaseUserExists(client *gophercloud.ServiceClient, dbmsID string, userName string, dbmsType string) (bool, *users.User, error) {
	result := users.Get(client, dbmsID, userName, dbmsType)
	if result.Err != nil {
		if errutil.IsNotFound(result.Err) {
			return false, nil, nil
		}
		return false, nil, result.Err
	}

	userObj, err := result.Extract()
	if err != nil {
		return false, nil, err
	}
	if userObj == nil {
		return false, nil, nil
	}

	return true, userObj, nil
}
