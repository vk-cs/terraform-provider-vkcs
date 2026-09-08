package flavors

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "flavors"
}

func flavorsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func flavorURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
