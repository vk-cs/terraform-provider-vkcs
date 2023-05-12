package networking_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccNetworkingSubnetRoute_importBasic(t *testing.T) {
	resourceName := "vkcs_networking_subnet_route.subnet_route_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
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
