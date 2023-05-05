package databases

import (
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db"
)

func baseURL(dbmsType string, dbmsID string) string {
	var pref string
	if dbmsType == db.DBMSTypeInstance {
		pref = "instances"
	} else {
		pref = "clusters"
	}
	return strings.Join([]string{pref, dbmsID, "databases"}, "/")
}

func databasesURL(c *gophercloud.ServiceClient, dbmsType string, dbmsID string) string {
	return c.ServiceURL(baseURL(dbmsType, dbmsID))
}

func databaseURL(c *gophercloud.ServiceClient, dbmsType string, dbmsID string, dbName string) string {
	return c.ServiceURL(baseURL(dbmsType, dbmsID), dbName)
}
