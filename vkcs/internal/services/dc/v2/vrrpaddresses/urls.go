package vrrpaddresses

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "direct_connect/dc_vrrp_addresses"
}

func vrrpAddressesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func vrrpAddressURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
