package keys

import (
	"net/url"

	"github.com/gophercloud/gophercloud"
)

func baseURL() string {
	return "transit"
}

func keysURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL(), "keys?list=true")
}

func keyURL(c *gophercloud.ServiceClient, keyName string) string {
	return c.ServiceURL(baseURL(), "keys", url.PathEscape(keyName))
}

func keyConfigURL(c *gophercloud.ServiceClient, keyName string) string {
	return c.ServiceURL(baseURL(), "keys", url.PathEscape(keyName), "config")
}

func encryptURL(c *gophercloud.ServiceClient, keyName string) string {
	return c.ServiceURL(baseURL(), "encrypt", url.PathEscape(keyName))
}

func decryptURL(c *gophercloud.ServiceClient, keyName string) string {
	return c.ServiceURL(baseURL(), "decrypt", url.PathEscape(keyName))
}
