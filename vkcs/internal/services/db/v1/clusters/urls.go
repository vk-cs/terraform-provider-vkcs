package clusters

import (
	"github.com/gophercloud/gophercloud"
)

func baseURL() string {
	return "clusters"
}

func clustersURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func clusterURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}

func capabilitiesURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id, "capabilities")
}

func backupScheduleURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id, "backup_schedule")

}
