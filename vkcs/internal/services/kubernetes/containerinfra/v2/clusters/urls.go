package clusters

import (
	"github.com/gophercloud/gophercloud"
)

// rootURL generates URL for clusters root
func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("clusters")
}

// resourceURL generates URL for cluster resource
func resourceURL(c *gophercloud.ServiceClient, clusterID string) string {
	return c.ServiceURL("clusters", clusterID)
}

// upgradeURL generates URL for cluster upgrading
func upgradeURL(c *gophercloud.ServiceClient, clusterID string) string {
	return resourceURL(c, clusterID) + ":upgrade"
}

// scaleURL generates URL for cluster scaling
func scaleURL(c *gophercloud.ServiceClient, clusterID string) string {
	return resourceURL(c, clusterID) + ":scaleControlPlane"
}

func kubeconfigURL(c *gophercloud.ServiceClient, clusterID string) string {
	return resourceURL(c, clusterID) + "/kube_config"
}
