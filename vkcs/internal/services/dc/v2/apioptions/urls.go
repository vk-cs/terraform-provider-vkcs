package apioptions

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "direct_connect/dc_api_options"
}

func apiOptionsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}
