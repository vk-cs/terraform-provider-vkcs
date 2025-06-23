package templates

import (
	"github.com/gophercloud/gophercloud"
)

func baseURL() string {
	return "templates"
}

func templatesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}
