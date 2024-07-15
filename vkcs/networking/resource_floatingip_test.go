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

func TestAccNetworkingFloatingIP_fixedip_bind(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckNetworkingFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingFloatingIPFixedIPBind1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &fip),
					testAccCheckNetworkingFloatingIPBoundToCorrectIP(&fip, "192.168.199.20"),
					resource.TestCheckResourceAttr("vkcs_networking_floatingip.fip_1", "description", "test"),
					resource.TestCheckResourceAttr("vkcs_networking_floatingip.fip_1", "fixed_ip", "192.168.199.20"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingFloatingIPFixedipBind2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &fip),
					testAccCheckNetworkingFloatingIPBoundToCorrectIP(&fip, "192.168.199.10"),
					resource.TestCheckResourceAttr("vkcs_networking_floatingip.fip_1", "description", ""),
					resource.TestCheckResourceAttr("vkcs_networking_floatingip.fip_1", "fixed_ip", "192.168.199.10"),
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

func testAccCheckNetworkingFloatingIPBoundToCorrectIP(fip *floatingips.FloatingIP, fixedIP string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if fip.FixedIP != fixedIP {
			return fmt.Errorf("Floating IP associated with wrong fixed ip")
		}

		return nil
	}
}

const testAccNetworkingFloatingIPBasic = `
resource "vkcs_networking_floatingip" "fip_1" {
  pool = "{{.ExtNetName}}"
  description = "test floating IP"
}
`

const testAccNetworkingFloatingIPFixedIPBind1 = `
{{.BaseExtNetwork}}

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id
}

resource "vkcs_networking_router_interface" "router_interface_1" {
  router_id = vkcs_networking_router.router_1.id
  subnet_id = vkcs_networking_subnet.subnet_1.id
}

resource "vkcs_networking_router" "router_1" {
  name = "router_1"
  external_network_id = data.vkcs_networking_network.extnet.id
}

resource "vkcs_networking_port" "port_1" {
  admin_state_up = "true"
  network_id = vkcs_networking_subnet.subnet_1.network_id

  fixed_ip {
    subnet_id = vkcs_networking_subnet.subnet_1.id
    ip_address = "192.168.199.10"
  }

  fixed_ip {
    subnet_id = vkcs_networking_subnet.subnet_1.id
    ip_address = "192.168.199.20"
  }
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = "{{.ExtNetName}}"
  description = "test"
  port_id = vkcs_networking_port.port_1.id
  fixed_ip = vkcs_networking_port.port_1.fixed_ip.1.ip_address
  depends_on = [
	vkcs_networking_router_interface.router_interface_1,
  ]
}
`

const testAccNetworkingFloatingIPFixedipBind2 = `
{{.BaseExtNetwork}}

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id
}

resource "vkcs_networking_router_interface" "router_interface_1" {
  router_id = vkcs_networking_router.router_1.id
  subnet_id = vkcs_networking_subnet.subnet_1.id
}

resource "vkcs_networking_router" "router_1" {
  name = "router_1"
  external_network_id = data.vkcs_networking_network.extnet.id
}

resource "vkcs_networking_port" "port_1" {
  admin_state_up = "true"
  network_id = vkcs_networking_subnet.subnet_1.network_id

  fixed_ip {
    subnet_id = vkcs_networking_subnet.subnet_1.id
    ip_address = "192.168.199.10"
  }

  fixed_ip {
    subnet_id = vkcs_networking_subnet.subnet_1.id
    ip_address = "192.168.199.20"
  }
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = "{{.ExtNetName}}"
  port_id = vkcs_networking_port.port_1.id
  fixed_ip = vkcs_networking_port.port_1.fixed_ip.0.ip_address
  depends_on = [
	vkcs_networking_router_interface.router_interface_1,
  ]
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
