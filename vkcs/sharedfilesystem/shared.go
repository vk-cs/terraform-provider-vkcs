package sharedfilesystem

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/sharedfilesystem/v2/shares"
)

const (
	SharedFilesystemMinMicroversion               = "2.7"
	sharedFilesystemSharedAccessCephXMicroversion = "2.13"
)

func getShareExportLocationPath(client *gophercloud.ServiceClient, id string) (string, error) {
	exportLocations, err := shares.ListExportLocations(client, id).Extract()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve export location from the Shares API: %s", err)
	}
	if len(exportLocations) > 0 {
		return exportLocations[0].Path, nil
	}
	return "", nil
}
