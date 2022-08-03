package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingRouter_importBasic(t *testing.T) {
	resourceName := "vkcs_networking_router.router_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingRouterBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
