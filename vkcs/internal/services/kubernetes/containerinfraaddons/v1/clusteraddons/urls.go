package clusteraddons

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "cluster_addons"
}

func clusterAddonURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
