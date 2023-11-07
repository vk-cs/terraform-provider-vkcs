package compute

import (
	"fmt"
	"log"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	iattachinterfaces "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/attachinterfaces"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func computeInterfaceAttachAttachFunc(
	computeClient *gophercloud.ServiceClient, instanceID, attachmentID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		va, err := iattachinterfaces.Get(computeClient, instanceID, attachmentID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
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

		va, err := iattachinterfaces.Get(computeClient, instanceID, attachmentID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return va, "DETACHED", nil
			}
			return va, "", err
		}

		err = iattachinterfaces.Delete(computeClient, instanceID, attachmentID).ExtractErr()
		if err != nil {
			if errutil.IsNotFound(err) {
				return va, "DETACHED", nil
			}

			if errutil.Is(err, 400) {
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
