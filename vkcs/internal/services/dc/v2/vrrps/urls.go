package vrrps

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "direct_connect/dc_vrrps"
}

func vrrpsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func vrrpURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
