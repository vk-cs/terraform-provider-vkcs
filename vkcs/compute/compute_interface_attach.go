package compute

import (
	"fmt"
	"log"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func computeInterfaceAttachAttachFunc(
	computeClient *gophercloud.ServiceClient, instanceID, attachmentID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		va, err := attachinterfaces.Get(computeClient, instanceID, attachmentID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "ATTACHING", nil
			}
			return va, "", err
		}

		return va, "ATTACHED", nil
	}
}

func computeInterfaceAttachDetachFunc(
	computeClient *gophercloud.ServiceClient, instanceID, attachmentID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to detach vkcs_compute_interface_attach %s from instance %s",
			attachmentID, instanceID)

		va, err := attachinterfaces.Get(computeClient, instanceID, attachmentID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "DETACHED", nil
			}
			return va, "", err
		}

		err = attachinterfaces.Delete(computeClient, instanceID, attachmentID).ExtractErr()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "DETACHED", nil
			}

			if _, ok := err.(gophercloud.ErrDefault400); ok {
				return nil, "", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] vkcs_compute_interface_attach %s is still active.", attachmentID)
		return nil, "", nil
	}
}

func ComputeInterfaceAttachParseID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("unable to determine vkcs_compute_interface_attach %s ID", id)
	}

	instanceID := idParts[0]
	attachmentID := idParts[1]

	return instanceID, attachmentID, nil
}
