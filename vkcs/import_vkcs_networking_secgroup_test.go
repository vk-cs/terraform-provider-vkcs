package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingSecGroup_importBasic(t *testing.T) {
	resourceName := "vkcs_networking_secgroup.secgroup_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
