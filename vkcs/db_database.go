package vkcs

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	dbmsTypeInstance = "instance"
	dbmsTypeCluster  = "cluster"
)

// Custom type implementation of gophercloud/database.DBPage
type DBPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks to see whether the collection is empty.
func (page DBPage) IsEmpty() (bool, error) {
	dbs, err := ExtractDBs(page)
	return len(dbs) == 0, err
}

// NextPageURL will retrieve the next page URL.
func (page DBPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractDBs will convert a generic pagination struct into a more
// relevant slice of DB structs.
func ExtractDBs(page pagination.Page) ([]databases.Database, error) {
	r := page.(DBPage)
	var s struct {
		Databases []databases.Database `json:"databases"`
	}
	err := r.ExtractInto(&s)
	return s.Databases, err
}

func databaseDatabaseStateRefreshFunc(client databaseClient, dbmsID string, databaseName string, dbmsType string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pages, err := databaseList(client, dbmsID, dbmsType).AllPages()
		if err != nil {
			return nil, "", fmt.Errorf("unable to retrieve vkcs database databases: %s", err)
		}

		allDatabases, err := ExtractDBs(pages)
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

func databaseDatabaseExists(client databaseClient, dbmsID string, databaseName string, dbmsType string) (bool, error) {
	var err error

	pages, err := databaseList(client, dbmsID, dbmsType).AllPages()
	if err != nil {
		return false, err
	}

	allDatabases, err := ExtractDBs(pages)
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
