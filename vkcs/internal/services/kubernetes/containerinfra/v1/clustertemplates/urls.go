package clustertemplates

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "clustertemplates"
}

func templatesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func templateURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
