package interfaces

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "direct_connect/dc_interfaces"
}

func interfacesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func interfaceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
