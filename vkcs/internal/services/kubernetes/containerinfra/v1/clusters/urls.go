package clusters

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "clusters"
}

func clustersURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func clusterURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}

func kubeConfigURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id, "kube_config")
}

func actionsURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id, "actions")
}

func upgradeURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id, "actions", "upgrade")
}
