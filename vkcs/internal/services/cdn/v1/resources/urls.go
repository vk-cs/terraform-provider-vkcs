package resources

import (
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud"
)

func baseURL(projectID string) string {
	return fmt.Sprintf("projects/%s/resources", projectID)
}

func resourcesURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL(baseURL(projectID))
}

func resourceURL(c *gophercloud.ServiceClient, projectID string, id int) string {
	return c.ServiceURL(baseURL(projectID), strconv.Itoa(id))
}

func shieldingURL(c *gophercloud.ServiceClient, projectID string, resourceID int) string {
	return c.ServiceURL(baseURL(projectID), strconv.Itoa(resourceID), "shielding")
}

func issueLetsEncryptURL(c *gophercloud.ServiceClient, projectID string, resourceID int) string {
	return c.ServiceURL(baseURL(projectID), strconv.Itoa(resourceID), "ssl", "le", "issue")
}

func prefetchContentURL(c *gophercloud.ServiceClient, projectID string, resourceID int) string {
	return c.ServiceURL(baseURL(projectID), strconv.Itoa(resourceID), "prefetch")
}
