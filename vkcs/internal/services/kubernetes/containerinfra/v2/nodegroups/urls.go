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

// kubeconfigURL retrieves volume types with corresponding azs
func volumeTypesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("available-resources/storage-classes")
}

// getByName retrieves node group by name and cluster id
func getByName(c *gophercloud.ServiceClient, clusterID, ngName string) string {
	return c.ServiceURL("clusters", clusterID, "node-groups", ngName)
}
