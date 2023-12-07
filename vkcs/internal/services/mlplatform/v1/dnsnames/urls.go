package dnsnames

import (
	"github.com/gophercloud/gophercloud"
)

func baseURL() string {
	return "dns_name"
}

func dnsNamesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}
