package datastores

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "datastores"
}

func datastoresURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func datastoreURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}

func datastoreParametersURL(c *gophercloud.ServiceClient, dsType string, dsVersion string) string {
	return c.ServiceURL(baseURL(), dsType, "versions", dsVersion, "parameters")
}

func datastoreCapabilitiesURL(c *gophercloud.ServiceClient, dsType string, dsVersion string) string {
	return c.ServiceURL(baseURL(), dsType, "versions", dsVersion, "capabilities")
}
