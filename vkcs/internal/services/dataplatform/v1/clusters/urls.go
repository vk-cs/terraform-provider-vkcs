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

func clusterSettingsURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id, "settings")
}

func clusterUsersURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id, "users")
}
