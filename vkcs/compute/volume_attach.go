package compute

import (
	"fmt"
	"log"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	ivolumes "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/blockstorage/v3/volumes"
	ivolumeattach "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/volumeattach"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func ComputeVolumeAttachParseID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("unable to determine vkcs_compute_volume_attach ID")
	}

	instanceID := parts[0]
	attachmentID := parts[1]

	return instanceID, attachmentID, nil
}

func computeVolumeAttachAttachFunc(computeClient *gophercloud.ServiceClient, blockStorageClient *gophercloud.ServiceClient, instanceID, attachmentID string, volumeID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		va, err := ivolumeattach.Get(computeClient, instanceID, attachmentID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return va, "ATTACHING", nil
			}
			return va, "", err
		}

		v, err := ivolumes.Get(blockStorageClient, volumeID).Extract()
		if err != nil {
			return va, "", err
		}
		if v.Status == "error" {
			return va, "", fmt.Errorf("volume entered unexpected error status")
		}
		if v.Status != "in-use" {
			return va, "ATTACHING", nil
		}

		return va, "ATTACHED", nil
	}
}

func computeVolumeAttachDetachFunc(computeClient *gophercloud.ServiceClient, instanceID, attachmentID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] vkcs_compute_volume_attach attempting to detach VKCS volume %s from instance %s",
			attachmentID, instanceID)

		va, err := ivolumeattach.Get(computeClient, instanceID, attachmentID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return va, "DETACHED", nil
			}
			return va, "", err
		}

		err = ivolumeattach.Delete(computeClient, instanceID, attachmentID).ExtractErr()
		if err != nil {
			if errutil.IsNotFound(err) {
				return va, "DETACHED", nil
			}

			if errutil.Is(err, 400) {
				return nil, "", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] vkcs_compute_volume_attach (%s/%s) is still active.", instanceID, attachmentID)
		return nil, "", nil
	}
}
