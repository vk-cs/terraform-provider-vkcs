package vkcs

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingSubnetRoute_basic(t *testing.T) {
	var (
		router  routers.Router
		network networks.Network
		subnet  subnets.Subnet
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSubnetRouteCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingRouterExists("vkcs_networking_router.router_1", &router),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingRouterInterfaceExists("vkcs_networking_router_interface.int_1"),
					testAccCheckNetworkingSubnetRouteExists("vkcs_networking_subnet_route.subnet_route_1"),
				),
			},
			{
				Config: testAccNetworkingSubnetRouteUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetRouteExists("vkcs_networking_subnet_route.subnet_route_1"),
					testAccCheckNetworkingSubnetRouteExists("vkcs_networking_subnet_route.subnet_route_2"),
				),
			},
			{
				Config: testAccNetworkingSubnetRouteDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetRouteEmpty("vkcs_networking_subnet.subnet_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSubnetRouteEmpty(n string) resource.TestCheckFunc {
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

		subnet, err := subnets.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if subnet.ID != rs.Primary.ID {
			return fmt.Errorf("Subnet not found")
		}

		if len(subnet.HostRoutes) != 0 {
			return fmt.Errorf("Invalid number of route entries: %d", len(subnet.HostRoutes))
		}

		return nil
	}
}

func testAccCheckNetworkingSubnetRouteExists(n string) resource.TestCheckFunc {
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

		subnet, err := subnets.Get(networkingClient, rs.Primary.Attributes["subnet_id"]).Extract()
		if err != nil {
			return err
		}

		if subnet.ID != rs.Primary.Attributes["subnet_id"] {
			return fmt.Errorf("Subnet for route not found")
		}

		var found = false
		for _, r := range subnet.HostRoutes {
			if r.DestinationCIDR == rs.Primary.Attributes["destination_cidr"] && r.NextHop == rs.Primary.Attributes["next_hop"] {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Could not find route for destination CIDR: %s, next hop: %s", rs.Primary.Attributes["destination_cidr"], rs.Primary.Attributes["next_hop"])
		}

		return nil
	}
}

func testAccCheckNetworkingSubnetRouteDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	networkingClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_networking_subnet_route" {
			continue
		}

		var routeExists = false

		subnet, err := subnets.Get(networkingClient, rs.Primary.Attributes["subnet_id"]).Extract()
		if err == nil {
			var rts = subnet.HostRoutes
			for _, r := range rts {
				if r.DestinationCIDR == rs.Primary.Attributes["destination_cidr"] && r.NextHop == rs.Primary.Attributes["next_hop"] {
					routeExists = true
					break
				}
			}
		}

		if routeExists {
			return fmt.Errorf("Route still exists")
		}
	}

	return nil
}

const testAccNetworkingSubnetRouteCreate = `
resource "vkcs_networking_router" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.1"
  }
}

resource "vkcs_networking_router_interface" "int_1" {
  router_id = "${vkcs_networking_router.router_1.id}"
  port_id = "${vkcs_networking_port.port_1.id}"
}

resource "vkcs_networking_subnet_route" "subnet_route_1" {
  destination_cidr = "10.0.1.0/24"
  next_hop = "192.168.199.254"

  subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}
`

const testAccNetworkingSubnetRouteUpdate = `
resource "vkcs_networking_router" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.1"
  }
}

resource "vkcs_networking_router_interface" "int_1" {
  router_id = "${vkcs_networking_router.router_1.id}"
  port_id = "${vkcs_networking_port.port_1.id}"
}

resource "vkcs_networking_subnet_route" "subnet_route_1" {
  destination_cidr = "10.0.1.0/24"
  next_hop = "192.168.199.254"

  subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}

resource "vkcs_networking_subnet_route" "subnet_route_2" {
  destination_cidr = "10.0.2.0/24"
  next_hop = "192.168.199.254"

  subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}
`

const testAccNetworkingSubnetRouteDestroy = `
resource "vkcs_networking_router" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.1"
  }
}

resource "vkcs_networking_router_interface" "int_1" {
  router_id = "${vkcs_networking_router.router_1.id}"
  port_id = "${vkcs_networking_port.port_1.id}"
}
`
