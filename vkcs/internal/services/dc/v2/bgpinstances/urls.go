package bgpinstances

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "direct_connect/dc_bgps"
}

func bgpsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func bgpURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
