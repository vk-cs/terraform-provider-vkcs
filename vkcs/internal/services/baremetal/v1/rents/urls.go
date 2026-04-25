package rents

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "rent-requests"
}

func rentRequestsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func rentRequestURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
