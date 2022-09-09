package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingNetworkDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkDataSourceNetwork,
			},
			{
				Config: testAccNetworkingNetworkDataSourceBasic(),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkDataSourceNetwork,
			},
			{
				Config: testAccNetworkingNetworkDataSourceSubnet(),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingNetworkDataSourceNetwork,
			},
			{
				Config: testAccNetworkingNetworkDataSourceNetworkID(),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccNetworkingNetworkDataSourceExternalExplicit, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkDataSourceID("data.vkcs_networking_network.network_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "name", osExtNetName),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "external", "true"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "all_tags.#", "0"),
				),
			},
		},
	})
}

func TestAccNetworkingNetworkDataSource_externalImplicit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccNetworkingNetworkDataSourceExternalImplicit, testAccValues),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkDataSourceID("data.vkcs_networking_network.network_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "name", osExtNetName),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "external", "true"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_network.network_1", "all_tags.#", "0"),
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

func testAccNetworkingNetworkDataSourceBasic() string {
	return fmt.Sprintf(`
%s

data "vkcs_networking_network" "network_1" {
  name = vkcs_networking_network.network_1.name
  description = vkcs_networking_network.network_1.description
}
`, testAccNetworkingNetworkDataSourceNetwork)
}

func testAccNetworkingNetworkDataSourceSubnet() string {
	return fmt.Sprintf(`
%s

data "vkcs_networking_network" "network_1" {
  matching_subnet_cidr = vkcs_networking_subnet.subnet_1.cidr
  tags = [
    "foo",
    "bar",
  ]
}
`, testAccNetworkingNetworkDataSourceNetwork)
}

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

func testAccNetworkingNetworkDataSourceNetworkID() string {
	return fmt.Sprintf(`
%s

data "vkcs_networking_network" "network_1" {
  network_id = vkcs_networking_network.network_1.id
}
`, testAccNetworkingNetworkDataSourceNetwork)
}
