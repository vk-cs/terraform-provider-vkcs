package nodegroups

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "nodegroups"
}

func nodeGroupsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func nodeGroupURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}

func scaleURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id, "actions", "scale")
}
