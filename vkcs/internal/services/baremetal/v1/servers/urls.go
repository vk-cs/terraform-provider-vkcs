package servers

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "servers"
}

func serversURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func serverURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}

func provisionURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id, "provision")
}
