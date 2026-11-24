package restores

import "github.com/gophercloud/gophercloud"

func restoresURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("restores")
}

func restoreURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("restores", id)
}
