package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccNetworkingSubnetDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSubnetDataSourceSubnet,
			},
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingSubnetDataSourceBasic, map[string]string{"TestAccNetworkingSubnetDataSourceSubnet": testAccNetworkingSubnetDataSourceSubnet}),
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

func TestAccNetworkingSubnetDataSource_migrateToFramework(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"vkcs": {
						VersionConstraint: "0.3.0",
						Source:            "vk-cs/vkcs",
					},
				},
				Config: acctest.AccTestRenderConfig(testAccNetworkingSubnetDataSourceBasic, map[string]string{"TestAccNetworkingSubnetDataSourceSubnet": testAccNetworkingSubnetDataSourceSubnet}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetDataSourceID("data.vkcs_networking_subnet.subnet_1"),
					testAccCheckNetworkingSubnetDataSourceGoodNetwork("data.vkcs_networking_subnet.subnet_1", "vkcs_networking_network.network_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "name", "subnet_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "all_tags.#", "2"),
				),
			},
			{
				ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
				Config:                   acctest.AccTestRenderConfig(testAccNetworkingSubnetDataSourceBasic, map[string]string{"TestAccNetworkingSubnetDataSourceSubnet": testAccNetworkingSubnetDataSourceSubnet}),
				PlanOnly:                 true,
			},
		},
	})
}

func TestAccNetworkingSubnetDataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSubnetDataSourceSubnet,
			},
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingSubnetDataSourceCidr, map[string]string{"TestAccNetworkingSubnetDataSourceSubnet": testAccNetworkingSubnetDataSourceSubnet}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetDataSourceID("data.vkcs_networking_subnet.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "description", "my subnet description"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "all_tags.#", "2"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingSubnetDataSourceDhcpEnabled, map[string]string{"TestAccNetworkingSubnetDataSourceSubnet": testAccNetworkingSubnetDataSourceSubnet}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetDataSourceID("data.vkcs_networking_subnet.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "all_tags.#", "2"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingSubnetDataSourceGatewayIP, map[string]string{"TestAccNetworkingSubnetDataSourceSubnet": testAccNetworkingSubnetDataSourceSubnet}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetDataSourceID("data.vkcs_networking_subnet.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_subnet.subnet_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

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

const testAccNetworkingSubnetDataSourceSubnet = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  description = "my subnet description"
  cidr = "192.168.199.0/24"
  network_id = vkcs_networking_network.network_1.id
  tags = [
    "foo",
    "bar",
  ]
}
`

const testAccNetworkingSubnetDataSourceBasic = `
{{.TestAccNetworkingSubnetDataSourceSubnet}}

data "vkcs_networking_subnet" "subnet_1" {
  name = vkcs_networking_subnet.subnet_1.name
}
`

const testAccNetworkingSubnetDataSourceCidr = `
{{.TestAccNetworkingSubnetDataSourceSubnet}}

data "vkcs_networking_subnet" "subnet_1" {
  cidr = "192.168.199.0/24"
  tags = []
}
`

const testAccNetworkingSubnetDataSourceDhcpEnabled = `
{{.TestAccNetworkingSubnetDataSourceSubnet}}

data "vkcs_networking_subnet" "subnet_1" {
  network_id = vkcs_networking_network.network_1.id
  dhcp_enabled = true
  tags = [
    "bar",
  ]
}
`

const testAccNetworkingSubnetDataSourceGatewayIP = `
{{.TestAccNetworkingSubnetDataSourceSubnet}}

data "vkcs_networking_subnet" "subnet_1" {
  gateway_ip = vkcs_networking_subnet.subnet_1.gateway_ip
}
`
