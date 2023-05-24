package networking_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccNetworkingPort_importBasic(t *testing.T) {
	resourceName := "vkcs_networking_port.port_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"fixed_ip",
				},
			},
		},
	})
}

func TestAccNetworkingPort_importAllowedAddressPairs(t *testing.T) {
	resourceName := "vkcs_networking_port.instance_port"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortAllowedAddressPairs1,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"fixed_ip",
				},
			},
		},
	})
}

func TestAccNetworkingPort_importAllowedAddressPairsNoMAC(t *testing.T) {
	resourceName := "vkcs_networking_port.instance_port"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortAllowedAddressPairsNoMAC,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"fixed_ip",
				},
			},
		},
	})
}

func TestAccNetworkingPort_importDHCPOpts(t *testing.T) {
	resourceName := "vkcs_networking_port.port_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortCreateExtraDhcpOpts,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"fixed_ip",
				},
			},
		},
	})
}
