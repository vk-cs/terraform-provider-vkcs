package blockstorage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccBlockStorageSnapshot_importBasic(t *testing.T) {
	resourceName := "vkcs_blockstorage_snapshot.snapshot_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckBlockStorageSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageSnapshotBasic, map[string]string{"TestAccBlockStorageVolumeBasic": acctest.AccTestRenderConfig(testAccBlockStorageVolumeBasic)}),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
