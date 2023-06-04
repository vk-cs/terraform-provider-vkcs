package providers

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "providers"
}

func providersURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}
