package addons

import (
	"github.com/gophercloud/gophercloud"
)

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("addons")
}

func getAddonByNameAndVersion(c *gophercloud.ServiceClient, addonName, addonVersion string) string {
	return c.ServiceURL("addons", addonName, "versions", addonVersion)
}

func getAddonByGlobalID(c *gophercloud.ServiceClient, addonVersionID string) string {
	return c.ServiceURL("addon_versions", addonVersionID)
}

func createClusterAddon(c *gophercloud.ServiceClient, clusterID string) string {
	return c.ServiceURL("clusters", clusterID, "addons")
}

func getClusterAddon(c *gophercloud.ServiceClient, clusterAddonID string) string {
	return c.ServiceURL("cluster_addons", clusterAddonID)
}

func getClusterAddonByClusterAndName(c *gophercloud.ServiceClient, clusterID, baseAddonName string) string {
	return c.ServiceURL("clusters", clusterID, "addons", baseAddonName)
}

func updateClusterAddon(c *gophercloud.ServiceClient, clusterID, clusterAddonID string) string {
	return c.ServiceURL("clusters", clusterID, "addons", clusterAddonID)
}

func deleteClusterAddon(c *gophercloud.ServiceClient, clusterID, clusterAddonID string) string {
	return c.ServiceURL("clusters", clusterID, "addons", clusterAddonID)
}
