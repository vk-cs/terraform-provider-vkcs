package vkcs

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type volumeChangeTypeOpts struct {
	volumeactions.ChangeTypeOpts
	AvailabilityZone string `json:"availability_zone,omitempty"`
}

func (opts volumeChangeTypeOpts) ToVolumeChangeTypeMap() (map[string]interface{}, error) {
	return BuildRequest(opts, "os-retype")
}

func blockStorageVolumeStateRefreshFunc(client *gophercloud.ServiceClient, volumeID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := volumes.Get(client, volumeID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return v, bsVolumeStatusDeleted, nil
			}
			return nil, "", err
		}
		if v.Status == "error" {
			return v, v.Status, fmt.Errorf("The volume is in error status. " +
				"Please check with VKCS support or check the Block Storage " +
				"API logs to see why this error occurred.")
		}

		return v, v.Status, nil
	}
}
