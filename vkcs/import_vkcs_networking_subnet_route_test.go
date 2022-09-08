package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingSubnetRoute_importBasic(t *testing.T) {
	resourceName := "vkcs_networking_subnet_route.subnet_route_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSubnetRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSubnetRouteCreate,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
