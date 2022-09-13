package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

func TestAccNetworkingSubnet_basic(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSubnetBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingSubnetDNSConsistency("vkcs_networking_subnet.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "allocation_pool.0.start", "192.168.199.100"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "description", "my subnet description"),
				),
			},
			{
				Config: testAccNetworkingSubnetUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "name", "subnet_1"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "gateway_ip", "192.168.199.1"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "enable_dhcp", "true"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "allocation_pool.0.start", "192.168.199.150"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "description", ""),
				),
			},
		},
	})
}

func TestAccNetworkingSubnet_enableDHCP(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSubnetEnableDhcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "enable_dhcp", "true"),
				),
			},
		},
	})
}

func TestAccNetworkingSubnet_disableDHCP(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSubnetDisableDhcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "enable_dhcp", "false"),
				),
			},
		},
	})
}

func TestAccNetworkingSubnet_noGateway(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSubnetNoGateway,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "gateway_ip", ""),
				),
			},
		},
	})
}

func TestAccNetworkingSubnet_impliedGateway(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSubnetImpliedGateway,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "gateway_ip", "192.168.199.1"),
				),
			},
		},
	})
}

func TestAccNetworkingSubnet_timeout(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSubnetTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
				),
			},
		},
	})
}

func TestAccNetworkingSubnet_multipleAllocationPools(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSubnetMultipleAllocationPools1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "allocation_pool.#", "2"),
				),
			},
			{
				Config: testAccNetworkingSubnetMultipleAllocationPools2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "allocation_pool.#", "2"),
				),
			},
			{
				Config: testAccNetworkingSubnetMultipleAllocationPools3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "allocation_pool.#", "2"),
				),
			},
		},
	})
}

func TestAccNetworkingSubnet_clearDNSNameservers(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSubnetClearDNSNameservers1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingSubnetDNSConsistency("vkcs_networking_subnet.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "dns_nameservers.#", "2"),
				),
			},
			{
				Config: testAccNetworkingSubnetClearDNSNameservers2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_networking_subnet.subnet_1", "dns_nameservers.#", "0"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSubnetDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	networkingClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_networking_subnet" {
			continue
		}

		_, err := subnets.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingSubnetExists(n string, subnet *subnets.Subnet) resource.TestCheckFunc {
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

		found, err := subnets.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Subnet not found")
		}

		*subnet = *found

		return nil
	}
}

func testAccCheckNetworkingSubnetDNSConsistency(n string, subnet *subnets.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		for i, dns := range subnet.DNSNameservers {
			if dns != rs.Primary.Attributes[fmt.Sprintf("dns_nameservers.%d", i)] {
				return fmt.Errorf("Dns Nameservers list elements or order is not consistent")
			}
		}

		return nil
	}
}

const testAccNetworkingSubnetBasic = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  description = "my subnet description"
  cidr = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id

  dns_nameservers = ["10.0.16.4", "213.186.33.99"]

  allocation_pool {
    start = "192.168.199.100"
    end = "192.168.199.200"
  }
}
`

const testAccNetworkingSubnetUpdate = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  gateway_ip = "192.168.199.1"
  network_id = vkcs_networking_network.network_1.id

  dns_nameservers = ["10.0.16.4", "213.186.33.99"]

  allocation_pool {
    start = "192.168.199.150"
    end = "192.168.199.200"
  }
}
`

const testAccNetworkingSubnetEnableDhcp = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  gateway_ip = "192.168.199.1"
  enable_dhcp = true
  network_id = vkcs_networking_network.network_1.id
}
`

const testAccNetworkingSubnetDisableDhcp = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  enable_dhcp = false
  network_id = vkcs_networking_network.network_1.id
}
`

const testAccNetworkingSubnetNoGateway = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  no_gateway = true
  network_id = vkcs_networking_network.network_1.id
}
`

const testAccNetworkingSubnetImpliedGateway = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}
resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id
}
`

const testAccNetworkingSubnetTimeout = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  cidr = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id

  allocation_pool {
    start = "192.168.199.100"
    end = "192.168.199.200"
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingSubnetMultipleAllocationPools1 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "10.3.0.0/16"
  network_id = vkcs_networking_network.network_1.id

  allocation_pool {
    start = "10.3.0.2"
    end = "10.3.0.255"
  }

  allocation_pool {
    start = "10.3.255.0"
    end = "10.3.255.254"
  }
}
`

const testAccNetworkingSubnetMultipleAllocationPools2 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "10.3.0.0/16"
  network_id = vkcs_networking_network.network_1.id

  allocation_pool {
    start = "10.3.255.0"
    end = "10.3.255.254"
  }

  allocation_pool {
    start = "10.3.0.2"
    end = "10.3.0.255"
  }
}
`

const testAccNetworkingSubnetMultipleAllocationPools3 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "10.3.0.0/16"
  network_id = vkcs_networking_network.network_1.id

  allocation_pool {
    start = "10.3.255.10"
    end = "10.3.255.154"
  }

  allocation_pool {
    start = "10.3.0.2"
    end = "10.3.0.255"
  }
}
`

const testAccNetworkingSubnetClearDNSNameservers1 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  cidr = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id

  dns_nameservers = ["10.0.16.4", "213.186.33.99"]

  allocation_pool {
    start = "192.168.199.100"
    end = "192.168.199.200"
  }
}
`

const testAccNetworkingSubnetClearDNSNameservers2 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id

  allocation_pool {
    start = "192.168.199.100"
    end = "192.168.199.200"
  }
}
`
