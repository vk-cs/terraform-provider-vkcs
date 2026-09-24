package secrets

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "secret"
}

func keysURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL(), "metadata")
}

func keyURL(c *gophercloud.ServiceClient, keyName string) string {
	return c.ServiceURL(baseURL(), "data", keyName)
}
