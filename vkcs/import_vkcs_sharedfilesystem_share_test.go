package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSFSShare_importBasic(t *testing.T) {
	resourceName := "vkcs_sharedfilesystem_share.share_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckSFS(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSShareConfigBasic(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
