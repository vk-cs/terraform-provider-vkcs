package ipportforwardings

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "direct_connect/dc_ip_port_forwardings"
}

func ipPortForwardingsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func ipPortForwardingURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
