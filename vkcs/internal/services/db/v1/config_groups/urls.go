package configgroups

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "configurations"
}

func configGroupsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func configGroupURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
