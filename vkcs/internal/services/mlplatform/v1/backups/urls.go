package backups

import (
	"github.com/gophercloud/gophercloud"
)

func baseURL() string {
	return "backups"
}

func backupsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func backupURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
