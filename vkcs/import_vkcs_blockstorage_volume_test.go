package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockStorageVolume_importBasic(t *testing.T) {
	resourceName := "vkcs_blockstorage_volume.volume_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccBlockStorageVolumeBasic, testAccValues),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
