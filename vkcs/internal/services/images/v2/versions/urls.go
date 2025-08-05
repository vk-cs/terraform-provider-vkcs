package versions

import (
	"strings"

	"github.com/gophercloud/gophercloud"
)

func versionsURL(c *gophercloud.ServiceClient) string {
	return strings.TrimSuffix(c.ResourceBaseURL(), "v2/")
}
