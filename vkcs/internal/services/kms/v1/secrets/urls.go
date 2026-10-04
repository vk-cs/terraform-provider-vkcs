package secrets

import (
	"net/url"

	"github.com/gophercloud/gophercloud"
)

func baseURL() string {
	return "secret"
}

func secretsMetadataURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL(), "metadata/?list=true")
}

func secretDataURL(c *gophercloud.ServiceClient, secretPath string) string {
	return c.ServiceURL(baseURL(), "data", url.PathEscape(secretPath))
}

func secretMetadataURL(c *gophercloud.ServiceClient, secretName string) string {
	return c.ServiceURL(baseURL(), "metadata", url.PathEscape(secretName))
}
