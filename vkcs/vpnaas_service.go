package vkcs

import (
	"github.com/gophercloud/gophercloud"
)

// Restart allows VPN services to be restarted.
func VPNAASServiceRestart(c *gophercloud.ServiceClient, id string) (r RestartResult) {
	resp, err := c.Get(vpnaasServiceRestartURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type RestartResult struct {
	gophercloud.ErrResult
}
