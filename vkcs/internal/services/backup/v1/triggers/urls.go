package triggers

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "triggers"
}

func triggerURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}

func triggersURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}
