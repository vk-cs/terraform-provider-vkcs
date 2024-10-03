package origingroups

import (
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud"
)

func baseURL(projectID string) string {
	return fmt.Sprintf("projects/%s/originGroups", projectID)
}

func originGroupsURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL(baseURL(projectID))
}

func originGroupURL(c *gophercloud.ServiceClient, projectID string, id int) string {
	return c.ServiceURL(baseURL(projectID), strconv.Itoa(id))
}
