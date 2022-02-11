package vkcs

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
)

func computeVolumeAttachParseID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("unable to determine vkcs_compute_volume_attach ID")
	}

	instanceID := parts[0]
	attachmentID := parts[1]

	return instanceID, attachmentID, nil
}

func computeVolumeAttachAttachFunc(computeClient *gophercloud.ServiceClient, blockStorageClient *gophercloud.ServiceClient, instanceID, attachmentID string, volumeID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		va, err := volumeattach.Get(computeClient, instanceID, attachmentID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "ATTACHING", nil
			}
			return va, "", err
		}

		v, err := volumes.Get(blockStorageClient, volumeID).Extract()
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

func computeVolumeAttachDetachFunc(computeClient *gophercloud.ServiceClient, instanceID, attachmentID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] vkcs_compute_volume_attach attempting to detach OpenStack volume %s from instance %s",
			attachmentID, instanceID)

		va, err := volumeattach.Get(computeClient, instanceID, attachmentID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "DETACHED", nil
			}
			return va, "", err
		}

		err = volumeattach.Delete(computeClient, instanceID, attachmentID).ExtractErr()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "DETACHED", nil
			}

			if _, ok := err.(gophercloud.ErrDefault400); ok {
				return nil, "", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] vkcs_compute_volume_attach (%s/%s) is still active.", instanceID, attachmentID)
		return nil, "", nil
	}
}
