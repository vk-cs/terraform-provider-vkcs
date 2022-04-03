package vkcs

import "github.com/gophercloud/gophercloud"

const (
	vpnaasRootPath    = "vpn"
	vpnaasServicePath = "vpnservices"
)

func vpnaasServiceRestartURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(vpnaasRootPath, vpnaasServicePath, id, "restart")
}
