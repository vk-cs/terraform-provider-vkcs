package templater

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
)

func baseMonitoringURL(projectID string) string {
	return fmt.Sprintf("project/%s", projectID)
}

func serviceUserURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL(baseMonitoringURL(projectID), "user")
}
