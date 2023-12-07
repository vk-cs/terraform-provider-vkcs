package instances

import (
	"github.com/gophercloud/gophercloud"
)

func baseURL() string {
	return "instances"
}

func instancesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func instanceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}

func instanceActionURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id, "action")
}
