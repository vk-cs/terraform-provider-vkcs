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

// kubeconfigURL generates URL for retrieving kubeconfig
func kubeconfigURL(c *gophercloud.ServiceClient, clusterID string) string {
	return resourceURL(c, clusterID) + "/kube_config"
}

// azURL generates URL for retrieving azs for clusters
func azURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("availability-zones")
}

// listKubeVersionURL generates URL for retrieving list of available k8s_versions
func listKubeVersionURL(c *gophercloud.ServiceClient) string {
	return rootURL(c) + "/k8s_versions"
}
