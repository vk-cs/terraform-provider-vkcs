package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingSubnetDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSubnetDataSourceSubnet,
			},
			{
				Config: testAccOpenStackNetworkingSubnetDataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetDataSourceID("data.vkcs_networking_subnet.subnet_1"),
					testAccCheckNetworkingSubnetDataSourceGoodNetwork("data.vkcs_networking_subnet.subnet_1", "vkcs_networking_network.network_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "name", "subnet_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

func TestAccNetworkingSubnetDataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSubnetDataSourceSubnet,
			},
			{
				Config: testAccOpenStackNetworkingSubnetDataSourceCidr(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetDataSourceID("data.vkcs_networking_subnet.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "description", "my subnet description"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "all_tags.#", "2"),
				),
			},
			{
				Config: testAccOpenStackNetworkingSubnetDataSourceDhcpEnabled(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetDataSourceID("data.vkcs_networking_subnet.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "all_tags.#", "2"),
				),
			},
			{
				Config: testAccOpenStackNetworkingSubnetDataSourceIPVersion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetDataSourceID("data.vkcs_networking_subnet.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "all_tags.#", "2"),
				),
			},
			{
				Config: testAccOpenStackNetworkingSubnetDataSourceGatewayIP(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetDataSourceID("data.vkcs_networking_subnet.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

// func TestAccNetworkingSubnetDataSource_networkIdAttribute(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckNetworking(t) },
// 		ProviderFactories: testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccOpenStackNetworkingSubnetDataSourceNetworkIDAttribute(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckNetworkingSubnetDataSourceID("data.vkcs_networking_subnet.subnet_1"),
// 					testAccCheckNetworkingSubnetDataSourceGoodNetwork("data.vkcs_networking_subnet.subnet_1", "vkcs_networking_network.network_1"),
// 					resource.TestCheckResourceAttr(
// 						"data.vkcs_networking_subnet.subnet_1", "tags.#", "1"),
// 					resource.TestCheckResourceAttr(
// 						"data.vkcs_networking_subnet.subnet_1", "all_tags.#", "2"),
// 					testAccCheckNetworkingPortV2ID("openstack_networking_port_v2.port_1"),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccNetworkingSubnetDataSource_subnetPoolIdAttribute(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckNetworking(t) },
// 		ProviderFactories: testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccOpenStackNetworkingSubnetDataSourceSubnetPoolIDAttribute(),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckNetworkingSubnetDataSourceID("data.vkcs_networking_subnet.subnet_1"),
// 					resource.TestCheckResourceAttr(
// 						"data.vkcs_networking_subnet.subnet_1", "tags.#", "2"),
// 					resource.TestCheckResourceAttr(
// 						"data.vkcs_networking_subnet.subnet_1", "all_tags.#", "2"),
// 				),
// 			},
// 		},
// 	})
// }

func testAccCheckNetworkingSubnetDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find subnet data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Subnet data source ID not set")
		}

		return nil
	}
}

func testAccCheckNetworkingSubnetDataSourceGoodNetwork(n1, n2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds1, ok := s.RootModule().Resources[n1]
		if !ok {
			return fmt.Errorf("Can't find subnet data source: %s", n1)
		}

		if ds1.Primary.ID == "" {
			return fmt.Errorf("Subnet data source ID not set")
		}

		rs2, ok := s.RootModule().Resources[n2]
		if !ok {
			return fmt.Errorf("Can't find network resource: %s", n2)
		}

		if rs2.Primary.ID == "" {
			return fmt.Errorf("Network resource ID not set")
		}

		if rs2.Primary.ID != ds1.Primary.Attributes["network_id"] {
			return fmt.Errorf("Network id and subnet network_id don't match")
		}

		return nil
	}
}

const testAccOpenStackNetworkingSubnetDataSourceSubnet = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  description = "my subnet description"
  cidr = "192.168.199.0/24"
  network_id = "${vkcs_networking_network.network_1.id}"
  tags = [
    "foo",
    "bar",
  ]
}
`

// const testAccOpenStackNetworkingSubnetDataSourceSubnetWithSubnetPool = `
// resource "vkcs_networking_network" "network_1" {
//   name = "network_1"
//   admin_state_up = "true"
// }

// resource "openstack_networking_subnetpool_v2" "subnetpool_1" {
//   name = "my_ipv4_pool"
//   prefixes = ["10.11.12.0/24"]
// }

// resource "vkcs_networking_subnet" "subnet_1" {
//   name = "subnet_1"
//   cidr = "10.11.12.0/25"
//   network_id = "${vkcs_networking_network.network_1.id}"
//   subnetpool_id = "${openstack_networking_subnetpool_v2.subnetpool_1.id}"
//   tags = [
//     "foo",
//     "bar",
//   ]
// }
// `

func testAccOpenStackNetworkingSubnetDataSourceBasic() string {
	return fmt.Sprintf(`
%s

data "vkcs_networking_subnet" "subnet_1" {
  name = "${vkcs_networking_subnet.subnet_1.name}"
}
`, testAccOpenStackNetworkingSubnetDataSourceSubnet)
}

func testAccOpenStackNetworkingSubnetDataSourceCidr() string {
	return fmt.Sprintf(`
%s

data "vkcs_networking_subnet" "subnet_1" {
  cidr = "192.168.199.0/24"
  tags = []
}
`, testAccOpenStackNetworkingSubnetDataSourceSubnet)
}

func testAccOpenStackNetworkingSubnetDataSourceDhcpEnabled() string {
	return fmt.Sprintf(`
%s

data "vkcs_networking_subnet" "subnet_1" {
  network_id = "${vkcs_networking_network.network_1.id}"
  dhcp_enabled = true
  tags = [
    "bar",
  ]
}
`, testAccOpenStackNetworkingSubnetDataSourceSubnet)
}

func testAccOpenStackNetworkingSubnetDataSourceIPVersion() string {
	return fmt.Sprintf(`
%s

data "vkcs_networking_subnet" "subnet_1" {
  network_id = "${vkcs_networking_network.network_1.id}"
  ip_version = 4
}
`, testAccOpenStackNetworkingSubnetDataSourceSubnet)
}

func testAccOpenStackNetworkingSubnetDataSourceGatewayIP() string {
	return fmt.Sprintf(`
%s

data "vkcs_networking_subnet" "subnet_1" {
  gateway_ip = "${vkcs_networking_subnet.subnet_1.gateway_ip}"
}
`, testAccOpenStackNetworkingSubnetDataSourceSubnet)
}

// func testAccOpenStackNetworkingSubnetDataSourceNetworkIDAttribute() string {
// 	return fmt.Sprintf(`
// %s

// data "vkcs_networking_subnet" "subnet_1" {
//   subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
//   tags = [
//     "foo",
//   ]
// }

// resource "openstack_networking_port_v2" "port_1" {
//   name            = "test_port"
//   network_id      = "${data.vkcs_networking_subnet.subnet_1.network_id}"
//   admin_state_up  = "true"
// }

// `, testAccOpenStackNetworkingSubnetDataSourceSubnet)
// }

// func testAccOpenStackNetworkingSubnetDataSourceSubnetPoolIDAttribute() string {
// 	return fmt.Sprintf(`
// %s

// data "vkcs_networking_subnet" "subnet_1" {
//   subnetpool_id = "${vkcs_networking_subnet.subnet_1.subnetpool_id}"
//   tags = [
//     "foo",
//     "bar",
//   ]
// }
// `, testAccOpenStackNetworkingSubnetDataSourceSubnetWithSubnetPool)
// }
