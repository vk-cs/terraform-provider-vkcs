package vkcs

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingRouterInterface_basic_subnet(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingRouterInterfaceBasicSubnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingRouterExists("vkcs_networking_router.router_1", &router),
					testAccCheckNetworkingRouterInterfaceExists("vkcs_networking_router_interface.int_1"),
				),
			},
		},
	})
}

// TestAccNetworkingRouterInterface_v6_subnet tests that multiple router interfaces for IPv6 subnets
// which are attached to the same port are handled properly.
func TestAccNetworkingRouterInterface_v6_subnet(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet1 subnets.Subnet
	var subnet2 subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingRouterInterfaceV6Subnet + testAccNetworkingRouterInterfaceV6SubnetSecondInterface,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet1),
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_2", &subnet2),
					testAccCheckNetworkingRouterExists("vkcs_networking_router.router_1", &router),
					testAccCheckNetworkingRouterInterfaceExists("vkcs_networking_router_interface.int_1"),
					testAccCheckNetworkingRouterInterfaceExists("vkcs_networking_router_interface.int_2"),
				),
			},
			{ // Make sure deleting one of the router interfaces does not remove the other one.
				Config: testAccNetworkingRouterInterfaceV6Subnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingRouterInterfaceExists("vkcs_networking_router_interface.int_1"),
				),
			},
		},
	})
}

func TestAccNetworkingRouterInterface_basic_port(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingRouterInterfaceBasicPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingRouterExists("vkcs_networking_router.router_1", &router),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckNetworkingRouterInterfaceExists("vkcs_networking_router_interface.int_1"),
				),
			},
		},
	})
}

func TestAccNetworkingRouterInterface_timeout(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingRouterInterfaceTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingRouterExists("vkcs_networking_router.router_1", &router),
					testAccCheckNetworkingRouterInterfaceExists("vkcs_networking_router_interface.int_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingRouterInterfaceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	networkingClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_networking_router_interface" {
			continue
		}

		_, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Router interface still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingRouterInterfaceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(configer)
		networkingClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		found, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Router interface not found")
		}

		return nil
	}
}

const testAccNetworkingRouterInterfaceBasicSubnet = `
resource "vkcs_networking_router" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "vkcs_networking_router_interface" "int_1" {
  subnet_id = vkcs_networking_subnet.subnet_1.id
  router_id = vkcs_networking_router.router_1.id
}

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = vkcs_networking_network.network_1.id
}
`

const testAccNetworkingRouterInterfaceV6Subnet = `
resource "vkcs_networking_router" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "vkcs_networking_router_interface" "int_1" {
  subnet_id = vkcs_networking_subnet.subnet_1.id
  router_id = vkcs_networking_router.router_1.id
}

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  cidr = "fd00:0:0:1::/64"
  ip_version = 6
  network_id = vkcs_networking_network.network_1.id
}

resource "vkcs_networking_subnet" "subnet_2" {
  cidr = "fd00:0:0:2::/64"
  ip_version = 6
  network_id = vkcs_networking_network.network_1.id
}
`

const testAccNetworkingRouterInterfaceV6SubnetSecondInterface = `
resource "vkcs_networking_router_interface" "int_2" {
  subnet_id = vkcs_networking_subnet.subnet_2.id
  router_id = vkcs_networking_router.router_1.id
}
`

const testAccNetworkingRouterInterfaceBasicPort = `
resource "vkcs_networking_router" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "vkcs_networking_router_interface" "int_1" {
  router_id = vkcs_networking_router.router_1.id
  port_id = vkcs_networking_port.port_1.id
}

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = vkcs_networking_network.network_1.id
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = vkcs_networking_network.network_1.id

  fixed_ip {
    subnet_id = vkcs_networking_subnet.subnet_1.id
    ip_address = "192.168.199.1"
  }
}
`

const testAccNetworkingRouterInterfaceTimeout = `
resource "vkcs_networking_router" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "vkcs_networking_router_interface" "int_1" {
  subnet_id = vkcs_networking_subnet.subnet_1.id
  router_id = vkcs_networking_router.router_1.id

  timeouts {
    create = "5m"
    delete = "5m"
  }
}

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = vkcs_networking_network.network_1.id
}
`
