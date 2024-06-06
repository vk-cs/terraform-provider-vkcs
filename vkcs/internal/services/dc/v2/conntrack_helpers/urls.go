package conntrackhelpers

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "direct_connect/dc_conntrack_helpers"
}

func conntrackHelpersURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func conntrackHelperURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
