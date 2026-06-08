package secpolicies

import (
	"github.com/gophercloud/gophercloud"
)

// rootURL generates URL for security policy root
func rootURL(c *gophercloud.ServiceClient, clusterID string) string {
	return c.ServiceURL("clusters", clusterID, "security_policy")
}

// getURL generates URL for security policy resource
func getURL(c *gophercloud.ServiceClient, clusterSecPolicyID string) string {
	return c.ServiceURL("cluster_security_policy", clusterSecPolicyID)
}

// resourceURL generates URL for security policy resource
func resourceURL(c *gophercloud.ServiceClient, clusterID, clusterSecPolicyID string) string {
	return c.ServiceURL("clusters", clusterID, "security_policy", clusterSecPolicyID)
}
