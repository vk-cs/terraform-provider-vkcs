package bgpneighbors

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "direct_connect/dc_bgp_neighbors"
}

func bgpNeighborsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func bgpNeighborURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
