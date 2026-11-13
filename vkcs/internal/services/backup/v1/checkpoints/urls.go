package checkpoints

import "github.com/gophercloud/gophercloud"

func checkpointsURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("checkpoints")
}

func checkpointURL(client *gophercloud.ServiceClient, providerID, checkpointID string) string {
	return client.ServiceURL("providers", providerID, "checkpoints", checkpointID)
}
