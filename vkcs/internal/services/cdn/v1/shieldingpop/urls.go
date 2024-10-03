package shieldingpop

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
)

func baseURL(projectID string) string {
	return fmt.Sprintf("projects/%s/shielding_pop", projectID)
}

func shieldingPopsURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL(baseURL(projectID))
}
