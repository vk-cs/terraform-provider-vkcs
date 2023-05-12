package networking_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccNetworkingRouterInterface_importBasic_port(t *testing.T) {
	resourceName := "vkcs_networking_router_interface.int_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckNetworkingRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingRouterInterfaceBasicPort,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingRouterInterface_importBasic_subnet(t *testing.T) {
	resourceName := "vkcs_networking_router_interface.int_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckNetworkingRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingRouterInterfaceBasicSubnet,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
