package networking_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccNetworkingRouterRoute_importBasic(t *testing.T) {
	resourceName := "vkcs_networking_router_route.router_route_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
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
