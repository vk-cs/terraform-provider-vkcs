package vkcs

import (
	"fmt"

	db "github.com/gophercloud/gophercloud/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/openstack/db/v1/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func extractDatabaseUserDatabases(v []interface{}) (db.BatchCreateOpts, error) {
	Batch := make(db.BatchCreateOpts, len(v))
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

func databaseUserStateRefreshFunc(client databaseClient, dbmsID string, userName string, dbmsType string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pages, err := userList(client, dbmsID, dbmsType).AllPages()
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

func databaseUserExists(client databaseClient, dbmsID string, userName string, dbmsType string) (bool, users.User, error) {
	var exists bool
	var err error
	var userObj users.User

	pages, err := userList(client, dbmsID, dbmsType).AllPages()
	if err != nil {
		return exists, userObj, err
	}

	allUsers, err := users.ExtractUsers(pages)
	if err != nil {
		return exists, userObj, err
	}

	for _, v := range allUsers {
		if v.Name == userName {
			exists = true
			return exists, v, nil
		}
	}

	return false, userObj, err
}
