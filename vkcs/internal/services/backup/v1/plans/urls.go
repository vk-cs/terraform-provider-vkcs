package plans

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "plans"
}

func planURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}

func plansURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}
