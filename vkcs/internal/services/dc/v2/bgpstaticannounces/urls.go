package bgpstaticannounces

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "direct_connect/dc_bgp_static_announces"
}

func bgpStaticAnnouncesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func bgpStaticAnnounceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
