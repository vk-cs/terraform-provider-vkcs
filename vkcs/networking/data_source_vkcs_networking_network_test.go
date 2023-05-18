package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccNetworkingNetworkDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkDataSourceNetwork,
			},
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingNetworkDataSourceBasic, map[string]string{"TestAccNetworkingNetworkDataSourceNetwork": testAccNetworkingNetworkDataSourceNetwork}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkDataSourceID("data.vkcs_networking_network.network_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "name", "tf_test_network"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "description", "my network description"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

func TestAccNetworkingNetworkDataSource_subnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkDataSourceNetwork,
			},
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingNetworkDataSourceSubnet, map[string]string{"TestAccNetworkingNetworkDataSourceNetwork": testAccNetworkingNetworkDataSourceNetwork}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkDataSourceID("data.vkcs_networking_network.network_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "name", "tf_test_network"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

func TestAccNetworkingNetworkDataSource_networkID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkDataSourceNetwork,
			},
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingNetworkDataSourceNetworkID, map[string]string{"TestAccNetworkingNetworkDataSourceNetwork": testAccNetworkingNetworkDataSourceNetwork}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkDataSourceID("data.vkcs_networking_network.network_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "name", "tf_test_network"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

func TestAccNetworkingNetworkDataSource_externalExplicit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingNetworkDataSourceExternalExplicit),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkDataSourceID("data.vkcs_networking_network.network_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "name", acctest.OsExtNetName),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "external", "true"),
				),
			},
		},
	})
}

func TestAccNetworkingNetworkDataSource_externalImplicit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingNetworkDataSourceExternalImplicit),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkDataSourceID("data.vkcs_networking_network.network_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "name", acctest.OsExtNetName),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "external", "true"),
				),
			},
		},
	})
}

func testAccCheckNetworkingNetworkDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find network data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Network data source ID not set")
		}

		return nil
	}
}

const testAccNetworkingNetworkDataSourceNetwork = `
resource "vkcs_networking_network" "network_1" {
  name = "tf_test_network"
  description = "my network description"
  admin_state_up = "true"
  tags = [
    "foo",
    "bar",
  ]
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "tf_test_subnet"
  cidr = "192.168.199.0/24"
  no_gateway = true
  network_id = vkcs_networking_network.network_1.id
}
`

const testAccNetworkingNetworkDataSourceBasic = `
{{.TestAccNetworkingNetworkDataSourceNetwork}}

data "vkcs_networking_network" "network_1" {
  name = vkcs_networking_network.network_1.name
  description = vkcs_networking_network.network_1.description
}
`

const testAccNetworkingNetworkDataSourceSubnet = `
	{{.TestAccNetworkingNetworkDataSourceNetwork}}

data "vkcs_networking_network" "network_1" {
  matching_subnet_cidr = vkcs_networking_subnet.subnet_1.cidr
  tags = [
    "foo",
    "bar",
  ]
}
`

const testAccNetworkingNetworkDataSourceExternalExplicit = `
data "vkcs_networking_network" "network_1" {
  name = "{{.ExtNetName}}"
  external = "true"
}
`

const testAccNetworkingNetworkDataSourceExternalImplicit = `
data "vkcs_networking_network" "network_1" {
  name = "{{.ExtNetName}}"
}
`

const testAccNetworkingNetworkDataSourceNetworkID = `
{{.TestAccNetworkingNetworkDataSourceNetwork}}

data "vkcs_networking_network" "network_1" {
  network_id = vkcs_networking_network.network_1.id
}
`
