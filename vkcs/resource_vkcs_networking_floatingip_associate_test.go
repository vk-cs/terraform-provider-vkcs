package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func TestAccNetworkingFloatingIPAssociate_basic(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingFloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingFloatingIPAssociateBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFloatingIPExists(
						"vkcs_networking_floatingip_associate.fip_1", &fip),
					resource.TestCheckResourceAttrPtr(
						"vkcs_networking_floatingip_associate.fip_1", "floating_ip", &fip.FloatingIP),
					resource.TestCheckResourceAttrPtr(
						"vkcs_networking_floatingip_associate.fip_1", "port_id", &fip.PortID),
				),
			},
		},
	})
}

func TestAccNetworkingFloatingIPAssociate_twoFixedIPs(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingFloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingFloatingIPAssociateTwoFixedIPs1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFloatingIPExists(
						"vkcs_networking_floatingip_associate.fip_1", &fip),
					resource.TestCheckResourceAttrPtr(
						"vkcs_networking_floatingip_associate.fip_1", "floating_ip", &fip.FloatingIP),
					resource.TestCheckResourceAttrPtr(
						"vkcs_networking_floatingip_associate.fip_1", "port_id", &fip.PortID),
					testAccCheckNetworkingFloatingIPBoundToCorrectIP(&fip, "192.168.199.20"),
					resource.TestCheckResourceAttr("vkcs_networking_floatingip_associate.fip_1", "fixed_ip", "192.168.199.20"),
				),
			},
			{
				Config: testAccNetworkingFloatingIPAssociateTwoFixedIPs2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFloatingIPExists(
						"vkcs_networking_floatingip_associate.fip_1", &fip),
					resource.TestCheckResourceAttrPtr(
						"vkcs_networking_floatingip_associate.fip_1", "floating_ip", &fip.FloatingIP),
					resource.TestCheckResourceAttrPtr(
						"vkcs_networking_floatingip_associate.fip_1", "port_id", &fip.PortID),
					testAccCheckNetworkingFloatingIPBoundToCorrectIP(&fip, "192.168.199.21"),
					resource.TestCheckResourceAttr("vkcs_networking_floatingip_associate.fip_1", "fixed_ip", "192.168.199.21"),
				),
			},
		},
	})
}

func testAccCheckNetworkingFloatingIPAssociateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	networkClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS network client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_networking_floatingip" {
			continue
		}

		fip, err := floatingips.Get(networkClient, rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return nil
			}

			return fmt.Errorf("Error retrieving Floating IP: %s", err)
		}

		if fip.PortID != "" {
			return fmt.Errorf("Floating IP is still associated")
		}
	}

	return nil
}

func testAccNetworkingFloatingIPAssociateBasic() string {
	return fmt.Sprintf(`
%s

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
  external_network_id = data.vkcs_networking_network.extnet.id
}

resource "vkcs_networking_port" "port_1" {
  admin_state_up = "true"
  network_id = "${vkcs_networking_subnet.subnet_1.network_id}"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.20"
  }
}

resource "vkcs_networking_floatingip" "fip_1" {
  depends_on = ["vkcs_networking_router.router_1"]
  pool = data.vkcs_networking_network.extnet.name
}

resource "vkcs_networking_floatingip_associate" "fip_1" {

  floating_ip = "${vkcs_networking_floatingip.fip_1.address}"
  port_id = "${vkcs_networking_port.port_1.id}"
}
`, testAccBaseExtNetwork)
}

func testAccNetworkingFloatingIPAssociateTwoFixedIPs1() string {
	return fmt.Sprintf(`
%s

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
  external_network_id = data.vkcs_networking_network.extnet.id
}

resource "vkcs_networking_port" "port_1" {
  admin_state_up = "true"
  network_id = "${vkcs_networking_subnet.subnet_1.network_id}"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.20"
  }

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.21"
  }
}

resource "vkcs_networking_floatingip" "fip_1" {
  depends_on = ["vkcs_networking_router.router_1"]
  pool = data.vkcs_networking_network.extnet.name
}

resource "vkcs_networking_floatingip_associate" "fip_1" {

  floating_ip = "${vkcs_networking_floatingip.fip_1.address}"
  port_id = "${vkcs_networking_port.port_1.id}"
  fixed_ip = "192.168.199.20"
}
`, testAccBaseExtNetwork)
}

func testAccNetworkingFloatingIPAssociateTwoFixedIPs2() string {
	return fmt.Sprintf(`
%s

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
  external_network_id = data.vkcs_networking_network.extnet.id
}

resource "vkcs_networking_port" "port_1" {
  admin_state_up = "true"
  network_id = "${vkcs_networking_subnet.subnet_1.network_id}"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.20"
  }

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.21"
  }
}

resource "vkcs_networking_floatingip" "fip_1" {
  depends_on = ["vkcs_networking_router.router_1"]
  pool = data.vkcs_networking_network.extnet.name
}

resource "vkcs_networking_floatingip_associate" "fip_1" {

  floating_ip = "${vkcs_networking_floatingip.fip_1.address}"
  port_id = "${vkcs_networking_port.port_1.id}"
  fixed_ip = "192.168.199.21"
}
`, testAccBaseExtNetwork)
}
