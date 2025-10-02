package serviceusers

import (
	"github.com/gophercloud/gophercloud"
)

func baseURL() string {
	return "service-users"
}

func serviceUsersURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func serviceUserURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
