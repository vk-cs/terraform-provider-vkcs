package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

type testNetworkWithExtensions struct {
	networks.Network
	portsecurity.PortSecurityExt
	policies.QoSPolicyExt
}

func TestAccNetworkingNetwork_basic(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "name", "network_1"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "description", "my network description"),
				),
			},
			{
				Config: testAccNetworkingNetworkUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "name", "network_2"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "description", ""),
				),
			},
		},
	})
}

func TestAccNetworkingNetwork_netstack(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkNetstack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingRouterExists("vkcs_networking_router.router_1", &router),
					testAccCheckNetworkingRouterInterfaceExists(
						"vkcs_networking_router_interface.ri_1"),
				),
			},
		},
	})
}

func TestAccNetworkingNetwork_fullstack(t *testing.T) {
	var instance servers.Server
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccNetworkingNetworkFullstack),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckComputeInstanceExists("vkcs_compute_instance.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccNetworkingNetwork_timeout(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
				),
			},
		},
	})
}

func TestAccNetworkingNetwork_adminStateUp_omit(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkAdminStateUpOmit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "admin_state_up", "true"),
					testAccCheckNetworkingNetworkAdminStateUp(&network, true),
				),
			},
		},
	})
}

func TestAccNetworkingNetwork_adminStateUp_true(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkAdminStateUpTrue,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "admin_state_up", "true"),
					testAccCheckNetworkingNetworkAdminStateUp(&network, true),
				),
			},
		},
	})
}

func TestAccNetworkingNetwork_adminStateUp_false(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkAdminStateUpFalse,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "admin_state_up", "false"),
					testAccCheckNetworkingNetworkAdminStateUp(&network, false),
				),
			},
		},
	})
}

func TestAccNetworkingNetwork_adminStateUp_update(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkAdminStateUpOmit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "admin_state_up", "true"),
					testAccCheckNetworkingNetworkAdminStateUp(&network, true),
				),
			},
			{
				Config: testAccNetworkingNetworkAdminStateUpFalse,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "admin_state_up", "false"),
					testAccCheckNetworkingNetworkAdminStateUp(&network, false),
				),
			},
		},
	})
}

func TestAccNetworkingNetwork_portSecurity_omit(t *testing.T) {
	var network testNetworkWithExtensions

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkAdminStateUpOmit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkWithExtensionsExists(
						"vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "port_security_enabled", "true"),
					testAccCheckNetworkingNetworkPortSecurityEnabled(&network, true),
				),
			},
			{
				Config: testAccNetworkingNetworkPortSecurityDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkWithExtensionsExists(
						"vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "port_security_enabled", "false"),
					testAccCheckNetworkingNetworkPortSecurityEnabled(&network, false),
				),
			},
			{
				Config: testAccNetworkingNetworkPortSecurityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkWithExtensionsExists(
						"vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "port_security_enabled", "true"),
					testAccCheckNetworkingNetworkPortSecurityEnabled(&network, true),
				),
			},
		},
	})
}

func TestAccNetworkingNetwork_portSecurity_disabled(t *testing.T) {
	var network testNetworkWithExtensions

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkPortSecurityDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkWithExtensionsExists(
						"vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "port_security_enabled", "false"),
					testAccCheckNetworkingNetworkPortSecurityEnabled(&network, false),
				),
			},
			{
				Config: testAccNetworkingNetworkPortSecurityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkWithExtensionsExists(
						"vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "port_security_enabled", "true"),
					testAccCheckNetworkingNetworkPortSecurityEnabled(&network, true),
				),
			},
		},
	})
}

func TestAccNetworkingNetwork_portSecurity_enabled(t *testing.T) {
	var network testNetworkWithExtensions

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkPortSecurityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkWithExtensionsExists(
						"vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "port_security_enabled", "true"),
					testAccCheckNetworkingNetworkPortSecurityEnabled(&network, true),
				),
			},
			{
				Config: testAccNetworkingNetworkPortSecurityDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkWithExtensionsExists(
						"vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "port_security_enabled", "false"),
					testAccCheckNetworkingNetworkPortSecurityEnabled(&network, false),
				),
			},
		},
	})
}

func TestAccNetworkingNetwork_privateDnsDomain(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccNetworkingNetworkPrivateDNSDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					resource.TestCheckResourceAttr(
						"vkcs_networking_network.network_1", "private_dns_domain", "test.domain."),
				),
			},
		},
	})
}

func testAccCheckNetworkingNetworkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(clients.Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName, networking.DefaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_networking_network" {
			continue
		}

		_, err := networks.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Network still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingNetworkExists(n string, network *networks.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(clients.Config)
		networkingClient, err := config.NetworkingV2Client(osRegionName, networking.DefaultSDN)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		found, err := networks.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Network not found")
		}

		*network = *found

		return nil
	}
}

func testAccCheckNetworkingNetworkWithExtensionsExists(n string, network *testNetworkWithExtensions) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(clients.Config)
		networkingClient, err := config.NetworkingV2Client(osRegionName, networking.DefaultSDN)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		var n testNetworkWithExtensions
		err = networks.Get(networkingClient, rs.Primary.ID).ExtractInto(&n)
		if err != nil {
			return err
		}

		if n.ID != rs.Primary.ID {
			return fmt.Errorf("Network not found")
		}

		*network = n

		return nil
	}
}

func testAccCheckNetworkingNetworkAdminStateUp(network *networks.Network, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if network.AdminStateUp != expected {
			return fmt.Errorf("Network has wrong admin_state_up. Expected %t, got %t", expected, network.AdminStateUp)
		}

		return nil
	}
}

func testAccCheckNetworkingNetworkPortSecurityEnabled(network *testNetworkWithExtensions, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if network.PortSecurityEnabled != expected {
			return fmt.Errorf("Network has wrong port_security_enabled. Expected %t, got %t", expected, network.PortSecurityEnabled)
		}

		return nil
	}
}

const testAccNetworkingNetworkBasic = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  description = "my network description"
  admin_state_up = "true"
}
`

const testAccNetworkingNetworkUpdate = `
resource "vkcs_networking_network" "network_1" {
  name = "network_2"
  admin_state_up = "true"
}
`

const testAccNetworkingNetworkNetstack = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.10.0/24"
  network_id = vkcs_networking_network.network_1.id
}

resource "vkcs_networking_router" "router_1" {
  name = "router_1"
}

resource "vkcs_networking_router_interface" "ri_1" {
  router_id = vkcs_networking_router.router_1.id
  subnet_id = vkcs_networking_subnet.subnet_1.id
}
`

const testAccNetworkingNetworkFullstack = `
{{.BaseImage}}
{{.BaseFlavor}}

resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  //security_group_ids = [vkcs_compute_secgroup.secgroup_1.id]
  network_id = vkcs_networking_network.network_1.id

  fixed_ip {
    subnet_id =  vkcs_networking_subnet.subnet_1.id
    ip_address =  "192.168.199.23"
  }
}

resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  //security_groups = [vkcs_compute_secgroup.secgroup_1.name]

  network {
    port = vkcs_networking_port.port_1.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccNetworkingNetworkTimeout = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingNetworkAdminStateUpOmit = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
}
`

const testAccNetworkingNetworkAdminStateUpTrue = `
resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}
`

const testAccNetworkingNetworkAdminStateUpFalse = `
resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "false"
}
`

const testAccNetworkingNetworkPortSecurityDisabled = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  port_security_enabled = "false"
}
`

const testAccNetworkingNetworkPortSecurityEnabled = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  port_security_enabled = "true"
}
`

const testAccNetworkingNetworkPrivateDNSDomain = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  description = "my network description"
  admin_state_up = "true"
  private_dns_domain = "test.domain."
}
`
