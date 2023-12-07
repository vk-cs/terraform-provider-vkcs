package mlplatform

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/mlplatform/v1/instances"
)

func TestFlattenVolumeOpts(t *testing.T) {
	volumes := []instances.VolumeResponse{
		{
			AvailabilityZone: "az1",
			CinderID:         "cid3",
			Name:             "data_volume2",
			Size:             70,
			VolumeType:       "volume_type",
		},
		{
			AvailabilityZone: "az1",
			CinderID:         "cid1",
			Name:             "boot_volume",
			Size:             50,
			VolumeType:       "volume_type",
		},
		{
			AvailabilityZone: "az1",
			CinderID:         "cid2",
			Name:             "data_volume1",
			Size:             60,
			VolumeType:       "volume_type",
		},
	}

	expectedBootVolume := &MLPlatformVolumeModel{
		Size:       types.Int64Value(50),
		VolumeType: types.StringValue("volume_type"),
		Name:       types.StringValue("boot_volume"),
		VolumeID:   types.StringValue("cid1"),
	}

	expectedDataVolumes := []*MLPlatformVolumeModel{
		{
			Size:       types.Int64Value(60),
			VolumeType: types.StringValue("volume_type"),
			Name:       types.StringValue("data_volume1"),
			VolumeID:   types.StringValue("cid2"),
		},
		{
			Size:       types.Int64Value(70),
			VolumeType: types.StringValue("volume_type"),
			Name:       types.StringValue("data_volume2"),
			VolumeID:   types.StringValue("cid3"),
		},
	}

	actualRoot, actualData, actualAZ := flattenVolumeOpts(volumes)
	assert.Equal(t, expectedBootVolume, actualRoot)
	assert.Equal(t, expectedDataVolumes, actualData)
	assert.Equal(t, "az1", actualAZ)
}
