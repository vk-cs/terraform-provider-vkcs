package checkpoints

import "github.com/gophercloud/gophercloud"

func checkpointsURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("checkpoints")
}

func checkpointURL(client *gophercloud.ServiceClient, providerID, checkpointID string) string {
	return client.ServiceURL("providers", providerID, "checkpoints", checkpointID)
}

func planCheckpointsURL(client *gophercloud.ServiceClient, providerID string) string {
	return client.ServiceURL("providers", providerID, "plan_checkpoints")
}
