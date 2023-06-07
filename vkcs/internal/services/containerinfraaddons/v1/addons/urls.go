package addons

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "addons"
}

func clusterAddonsURL(c *gophercloud.ServiceClient, clusterID string) string {
	return c.ServiceURL(baseURL(), "clusters", clusterID)
}

func clusterAvailableAddonsURL(c *gophercloud.ServiceClient, clusterID string) string {
	return c.ServiceURL(baseURL(), "clusters", clusterID, "available")
}

func clusterAvailableAddonURL(c *gophercloud.ServiceClient, clusterID, addonID string) string {
	return c.ServiceURL(baseURL(), "clusters", clusterID, "available", addonID)
}

func installAddonToClusterURL(c *gophercloud.ServiceClient, addonID, clusterID string) string {
	return c.ServiceURL(baseURL(), addonID, "clusters", clusterID)
}
