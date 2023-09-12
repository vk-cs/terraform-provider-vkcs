package routers

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "direct_connect/dc_routers"
}

func routersURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func routerURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
