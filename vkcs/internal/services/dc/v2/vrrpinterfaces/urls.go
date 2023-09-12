package vrrpinterfaces

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "direct_connect/dc_vrrp_interfaces"
}

func vrrpInterfacesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func vrrpInterfaceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
