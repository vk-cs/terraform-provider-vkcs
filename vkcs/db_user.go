package vkcs

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	db "github.com/gophercloud/gophercloud/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/openstack/db/v1/users"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Custom type implementation of gophercloud/users.UserPage
type DBUserPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks to see whether the collection is empty.
func (page DBUserPage) IsEmpty() (bool, error) {
	users, err := ExtractUsers(page)
	return len(users) == 0, err
}

// NextPageURL will retrieve the next page URL.
func (page DBUserPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractUsers will convert a generic pagination struct into a more
// relevant slice of User structs.
func ExtractUsers(r pagination.Page) ([]users.User, error) {
	var s struct {
		Users []users.User `json:"users"`
	}
	err := (r.(DBUserPage)).ExtractInto(&s)
	return s.Users, err
}

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

func databaseUserExists(client databaseClient, dbmsID string, userName string, dbmsType string) (bool, *users.User, error) {
	var err error

	pages, err := userList(client, dbmsID, dbmsType).AllPages()
	if err != nil {
		return false, nil, err
	}

	allUsers, err := ExtractUsers(pages)
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
