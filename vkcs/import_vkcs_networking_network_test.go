package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingNetwork_importBasic(t *testing.T) {
	resourceName := "vkcs_networking_network.network_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
