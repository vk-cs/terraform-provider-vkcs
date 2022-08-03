package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingRouterRoute_importBasic(t *testing.T) {
	resourceName := "vkcs_networking_router_route.router_route_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingRouterRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingRouterRouteCreate,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
