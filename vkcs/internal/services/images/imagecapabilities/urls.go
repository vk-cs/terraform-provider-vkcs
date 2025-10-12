package capabilities

import (
	"github.com/gophercloud/gophercloud"
)

func baseCapabilitiesURL() string {
	return "image_os_capabilities"
}

func imageCapabilitiesURL(c *gophercloud.ServiceClient, imageID string) string {
	return c.ServiceURL(baseCapabilitiesURL(), imageID)
}
