package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func TestAccNetworkingFloatingIP_basic(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckNetworkingFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingFloatingIPBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &fip),
					resource.TestCheckResourceAttr("vkcs_networking_floatingip.fip_1", "description", "test floating IP"),
				),
			},
		},
	})
}

func TestAccNetworkingFloatingIP_subnetIDs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckNetworkingFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingFloatingIPSubnetIDs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_networking_floatingip.fip_1", "description", "test"),
				),
			},
		},
	})
}

func TestAccNetworkingFloatingIP_timeout(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckNetworkingFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingFloatingIPTimeout),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &fip),
				),
			},
		},
	})
}

func testAccCheckNetworkingFloatingIPDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	networkClient, err := config.NetworkingV2Client(acctest.OsRegionName, networking.DefaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS floating IP: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_networking_floatingip" {
			continue
		}

		_, err := floatingips.Get(networkClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Floating IP still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingFloatingIPExists(n string, kp *floatingips.FloatingIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		networkClient, err := config.NetworkingV2Client(acctest.OsRegionName, networking.DefaultSDN)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		found, err := floatingips.Get(networkClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Floating IP not found")
		}

		*kp = *found

		return nil
	}
}

const testAccNetworkingFloatingIPBasic = `
resource "vkcs_networking_floatingip" "fip_1" {
  pool = "{{.ExtNetName}}"
  description = "test floating IP"
}
`

const testAccNetworkingFloatingIPTimeout = `
resource "vkcs_networking_floatingip" "fip_1" {
  pool = "{{.ExtNetName}}"
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingFloatingIPSubnetIDs = `
data "vkcs_networking_network" "ext_network" {
	name = "{{.ExtNetName}}"
}	

resource "vkcs_networking_floatingip" "fip_1" {
  pool = data.vkcs_networking_network.ext_network.name
  description = "test"
  subnet_ids = flatten([
    data.vkcs_networking_network.ext_network.id, # wrong UUID
    data.vkcs_networking_network.ext_network.subnets,
    data.vkcs_networking_network.ext_network.id, # wrong UUID again
  ])
}
`
