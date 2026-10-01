package secrets

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "secret"
}

func secretsMetadataURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL(), "metadata/?list=true")
}

func secretDataURL(c *gophercloud.ServiceClient, secretName string) string {
	return c.ServiceURL(baseURL(), "data", secretName)
}

func secretMetadataURL(c *gophercloud.ServiceClient, secretName string) string {
	return c.ServiceURL(baseURL(), "metadata", secretName)
}
