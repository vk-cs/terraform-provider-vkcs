package secrets

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "secret"
}

func secretsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL(), "metadata")
}

func secretURL(c *gophercloud.ServiceClient, secretName string) string {
	return c.ServiceURL(baseURL(), "data", secretName)
}
