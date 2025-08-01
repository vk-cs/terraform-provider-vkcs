package versions

import (
	"strings"

	"github.com/gophercloud/gophercloud"
)

func versionsURL(c *gophercloud.ServiceClient) string {
	baseURL := c.ResourceBaseURL()
	return strings.TrimSuffix(baseURL, "/v2/")
}
