package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockStorageSnapshot_importBasic(t *testing.T) {
	resourceName := "vkcs_blockstorage_snapshot.snapshot_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageSnapshotBasic(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
