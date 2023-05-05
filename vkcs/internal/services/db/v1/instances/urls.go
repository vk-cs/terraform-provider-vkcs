package instances

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "instances"
}

func instancesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func instanceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}

func actionURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id, "action")
}

func rootUserURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id, "root")
}

func capabilitiesURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id, "capabilities")
}

func backupScheduleURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id, "backup_schedule")
}
