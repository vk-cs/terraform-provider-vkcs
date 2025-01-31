package anycastips

import (
	"github.com/gophercloud/gophercloud"
)

func baseURL() string {
	return "anycastips"
}

func anycastIPsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func anycastIPURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
