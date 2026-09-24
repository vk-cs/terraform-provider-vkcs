package keys

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "transit"
}

func keysURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL(), "keys")
}

func keyURL(c *gophercloud.ServiceClient, keyName string) string {
	return c.ServiceURL(keysURL(c), keyName)
}
