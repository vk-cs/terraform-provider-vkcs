package users

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
	return strings.Join([]string{pref, dbmsID, "users"}, "/")
}

func usersURL(c *gophercloud.ServiceClient, dbmsType string, dbmsID string) string {
	return c.ServiceURL(baseURL(dbmsType, dbmsID))
}

func userURL(c *gophercloud.ServiceClient, dbmsType string, dbmsID string, userName string) string {
	return c.ServiceURL(baseURL(dbmsType, dbmsID), userName)
}

func userDatabasesURL(c *gophercloud.ServiceClient, dbmsType string, dbmsID string, userName string) string {
	return c.ServiceURL(baseURL(dbmsType, dbmsID), userName, "databases")
}

func userDatabaseURL(c *gophercloud.ServiceClient, dbmsType string, dbmsID string, userName string, databaseName string) string {
	return c.ServiceURL(baseURL(dbmsType, dbmsID), userName, "databases", databaseName)
}
