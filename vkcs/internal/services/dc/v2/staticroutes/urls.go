package staticroutes

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "direct_connect/dc_static_routes"
}

func staticRoutesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func staticRouteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
