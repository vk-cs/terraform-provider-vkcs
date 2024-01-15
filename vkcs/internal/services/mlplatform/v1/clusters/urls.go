package clusters

import (
	"github.com/gophercloud/gophercloud"
)

func baseURL() string {
	return "magnum"
}

func clustersURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL() + "/")
}

func clusterURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
