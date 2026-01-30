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

// upgradeURL updrades cluster's versions
func upgradeURL(c *gophercloud.ServiceClient, clusterID string) string {
	return resourceURL(c, clusterID) + ":upgrade"
}

// scaleURL scales cluster's master nodes to the new flavor
func scaleURL(c *gophercloud.ServiceClient, clusterID string) string {
	return resourceURL(c, clusterID) + ":scaleControlPlane"
}

// kubeconfigURL retrieves cluster's kubeconfig
func kubeconfigURL(c *gophercloud.ServiceClient, clusterID string) string {
	return resourceURL(c, clusterID) + "/kube_config"
}

// kubeconfigURL retrieves available Kubernetes versions
func listKubeVersionURL(c *gophercloud.ServiceClient) string {
	return rootURL(c) + "/k8s_versions"
}

// kubeconfigURL retrieves volume types with corresponding azs
func volumeTypesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("available-resources/storage-classes")
}
