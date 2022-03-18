package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func TestAccNetworkingFloatingIP_basic(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingFloatingIPBasic,
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
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingFloatingIPFixedIPBind1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &fip),
					testAccCheckNetworkingFloatingIPBoundToCorrectIP(&fip, "192.168.199.20"),
					resource.TestCheckResourceAttr("vkcs_networking_floatingip.fip_1", "description", "test"),
					resource.TestCheckResourceAttr("vkcs_networking_floatingip.fip_1", "fixed_ip", "192.168.199.20"),
				),
			},
			{
				Config: testAccNetworkingFloatingIPFixedipBind2(),
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
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingFloatingIPSubnetIDs(),
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
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingFloatingIPTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFloatingIPExists("vkcs_networking_floatingip.fip_1", &fip),
				),
			},
		},
	})
}

func testAccCheckNetworkingFloatingIPDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	networkClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack floating IP: %s", err)
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

		config := testAccProvider.Meta().(configer)
		networkClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
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
  description = "test floating IP"
}
`

func testAccNetworkingFloatingIPFixedIPBind1() string {
	return fmt.Sprintf(`
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_router_interface" "router_interface_1" {
  router_id = "${vkcs_networking_router.router_1.id}"
  subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}

resource "vkcs_networking_router" "router_1" {
  name = "router_1"
  external_network_id = "%s"
}

resource "vkcs_networking_port" "port_1" {
  admin_state_up = "true"
  network_id = "${vkcs_networking_subnet.subnet_1.network_id}"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.10"
  }

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.20"
  }
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = "%s"
  description = "test"
  port_id = "${vkcs_networking_port.port_1.id}"
  fixed_ip = "${vkcs_networking_port.port_1.fixed_ip.1.ip_address}"
}
`, osExtGwID, osPoolName)
}

func testAccNetworkingFloatingIPFixedipBind2() string {
	return fmt.Sprintf(`
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_router_interface" "router_interface_1" {
  router_id = "${vkcs_networking_router.router_1.id}"
  subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}

resource "vkcs_networking_router" "router_1" {
  name = "router_1"
  external_network_id = "%s"
}

resource "vkcs_networking_port" "port_1" {
  admin_state_up = "true"
  network_id = "${vkcs_networking_subnet.subnet_1.network_id}"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.10"
  }

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.20"
  }
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = "%s"
  port_id = "${vkcs_networking_port.port_1.id}"
  fixed_ip = "${vkcs_networking_port.port_1.fixed_ip.0.ip_address}"
}
`, osExtGwID, osPoolName)
}

const testAccNetworkingFloatingIPTimeout = `
resource "vkcs_networking_floatingip" "fip_1" {
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

func testAccNetworkingFloatingIPSubnetIDs() string {
	return fmt.Sprintf(`
data "vkcs_networking_network" "ext_network" {
  name = "%s"
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
`, osPoolName)
}
