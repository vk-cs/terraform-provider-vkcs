package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSFSShareNetwork_importBasic(t *testing.T) {
	resourceName := "vkcs_sharedfilesystem_sharenetwork.sharenetwork_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckSFS(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSShareNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSShareNetworkConfigBasic(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
