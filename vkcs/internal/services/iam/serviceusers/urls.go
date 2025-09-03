package serviceusers

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
)

func baseURL(projectID string) string {
	return fmt.Sprintf("tenants/%s/service-users", projectID)
}

func serviceUsersURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL(baseURL(projectID))
}

func serviceUserURL(c *gophercloud.ServiceClient, projectID, id string) string {
	return c.ServiceURL(baseURL(projectID), id)
}
