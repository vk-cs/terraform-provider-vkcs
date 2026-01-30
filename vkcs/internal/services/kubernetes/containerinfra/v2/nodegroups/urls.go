package nodegroups

import (
	"github.com/gophercloud/gophercloud"
)

// rootURL generates URL for nodegroups root
func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("node-groups")
}

// resourceURL generates URL for nodegroup resource
func resourceURL(c *gophercloud.ServiceClient, nodeGroupID string) string {
	return c.ServiceURL("node-groups", nodeGroupID)
}
