package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingFloatingIP_importBasic(t *testing.T) {
	resourceName := "vkcs_networking_floatingip.fip_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingFloatingIPBasic(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
