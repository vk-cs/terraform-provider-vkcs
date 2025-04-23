package blockstorage

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	ivolumes "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/blockstorage/v3/volumes"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type volumeChangeTypeOpts struct {
	volumeactions.ChangeTypeOpts
	AvailabilityZone string `json:"availability_zone,omitempty"`
}

func (opts volumeChangeTypeOpts) ToVolumeChangeTypeMap() (map[string]interface{}, error) {
	return util.BuildRequest(opts, "os-retype")
}

func BlockStorageVolumeStateRefreshFunc(client *gophercloud.ServiceClient, volumeID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := ivolumes.Get(client, volumeID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return v, BSVolumeStatusDeleted, nil
			}
			return nil, "", err
		}
		if v.Status == "error" {
			return v, v.Status, fmt.Errorf("the volume is in error status. Please check with your cloud admin or check the Block Storage API logs to see why this error occurred")
		}

		return v, v.Status, nil
	}
}
