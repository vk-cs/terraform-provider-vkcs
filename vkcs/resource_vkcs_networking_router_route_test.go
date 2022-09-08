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

func TestAccNetworkingRouterRoute_basic(t *testing.T) {
	var router routers.Router
	var network [2]networks.Network
	var subnet [2]subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingRouterRouteCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingRouterExists("vkcs_networking_router.router_1", &router),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network[0]),
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet[0]),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network[1]),
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet[1]),
					testAccCheckNetworkingRouterInterfaceExists("vkcs_networking_router_interface.int_1"),
					testAccCheckNetworkingRouterInterfaceExists("vkcs_networking_router_interface.int_2"),
					testAccCheckNetworkingRouterRouteExists("vkcs_networking_router_route.router_route_1"),
				),
			},
			{
				Config: testAccNetworkingRouterRouteUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingRouterRouteExists("vkcs_networking_router_route.router_route_1"),
					testAccCheckNetworkingRouterRouteExists("vkcs_networking_router_route.router_route_2"),
				),
			},
			{
				Config: testAccNetworkingRouterRouteDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingRouterRouteEmpty("vkcs_networking_router.router_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingRouterRouteEmpty(n string) resource.TestCheckFunc {
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

		router, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if router.ID != rs.Primary.ID {
			return fmt.Errorf("Router not found")
		}

		if len(router.Routes) != 0 {
			return fmt.Errorf("Invalid number of route entries: %d", len(router.Routes))
		}

		return nil
	}
}

func testAccCheckNetworkingRouterRouteExists(n string) resource.TestCheckFunc {
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

		router, err := routers.Get(networkingClient, rs.Primary.Attributes["router_id"]).Extract()
		if err != nil {
			return err
		}

		if router.ID != rs.Primary.Attributes["router_id"] {
			return fmt.Errorf("Router for route not found")
		}

		found := false
		for _, r := range router.Routes {
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

func testAccCheckNetworkingRouterRouteDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	networkingClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_networking_router_route" {
			continue
		}

		var routeExists = false

		router, err := routers.Get(networkingClient, rs.Primary.Attributes["router_id"]).Extract()
		if err == nil {
			var rts = router.Routes
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

const testAccNetworkingRouterRouteCreate = `
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

resource "vkcs_networking_network" "network_2" {
  name = "network_2"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_2" {
  cidr = "192.168.200.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_2.id}"
}

resource "vkcs_networking_port" "port_2" {
  name = "port_2"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_2.id}"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_2.id}"
    ip_address = "192.168.200.1"
  }
}

resource "vkcs_networking_router_interface" "int_2" {
  router_id = "${vkcs_networking_router.router_1.id}"
  port_id = "${vkcs_networking_port.port_2.id}"
  depends_on = [vkcs_networking_router_interface.int_1]
}

resource "vkcs_networking_router_route" "router_route_1" {
  destination_cidr = "10.0.1.0/24"
  next_hop = "192.168.199.254"

  depends_on = [vkcs_networking_router_interface.int_1, vkcs_networking_router_interface.int_2]
  router_id = "${vkcs_networking_router.router_1.id}"
}
`

const testAccNetworkingRouterRouteUpdate = `
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

resource "vkcs_networking_network" "network_2" {
  name = "network_2"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_2" {
  cidr = "192.168.200.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_2.id}"
}

resource "vkcs_networking_port" "port_2" {
  name = "port_2"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_2.id}"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_2.id}"
    ip_address = "192.168.200.1"
  }
}

resource "vkcs_networking_router_interface" "int_2" {
  router_id = "${vkcs_networking_router.router_1.id}"
  port_id = "${vkcs_networking_port.port_2.id}"
  depends_on = [vkcs_networking_router_interface.int_1]
}

resource "vkcs_networking_router_route" "router_route_1" {
  destination_cidr = "10.0.1.0/24"
  next_hop = "192.168.199.254"

  depends_on = [vkcs_networking_router_interface.int_1, vkcs_networking_router_interface.int_2]
  router_id = "${vkcs_networking_router.router_1.id}"
}

resource "vkcs_networking_router_route" "router_route_2" {
  destination_cidr = "10.0.2.0/24"
  next_hop = "192.168.200.254"

  depends_on = [vkcs_networking_router_interface.int_1, vkcs_networking_router_interface.int_2]
  router_id = "${vkcs_networking_router.router_1.id}"
}
`

const testAccNetworkingRouterRouteDestroy = `
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

resource "vkcs_networking_network" "network_2" {
  name = "network_2"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_2" {
  ip_version = 4
  cidr = "192.168.200.0/24"
  network_id = "${vkcs_networking_network.network_2.id}"
}

resource "vkcs_networking_port" "port_2" {
  name = "port_2"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_2.id}"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_2.id}"
    ip_address = "192.168.200.1"
  }
}

resource "vkcs_networking_router_interface" "int_2" {
  router_id = "${vkcs_networking_router.router_1.id}"
  port_id = "${vkcs_networking_port.port_2.id}"
}
`
