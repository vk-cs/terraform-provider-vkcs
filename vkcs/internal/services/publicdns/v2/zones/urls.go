package zones

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "dns"
}

func zoneURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}

func zonesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL(), "")
}
