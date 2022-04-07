package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingRouterDataSource_name(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingRouterDataSourceRouter,
			},
			{
				Config: testAccNetworkingRouterDataSourceName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingRouterDataSourceID("data.vkcs_networking_router.router"),
					resource.TestCheckResourceAttrSet(
						"data.vkcs_networking_router.router", "name"),
					resource.TestCheckResourceAttrSet(
						"data.vkcs_networking_router.router", "description"),
					resource.TestCheckResourceAttrSet(
						"data.vkcs_networking_router.router", "admin_state_up"),
					resource.TestCheckResourceAttrSet(
						"data.vkcs_networking_router.router", "status"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_router.router", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_router.router", "all_tags.#", "2"),
				),
			},
		},
	})
}

func testAccCheckNetworkingRouterDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find router data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Router data source ID not set")
		}

		return nil
	}
}

const testAccNetworkingRouterDataSourceRouter = `
resource "vkcs_networking_router" "router" {
  name           = "router_tf"
  description    = "description"
  admin_state_up = "true"
  tags = [
    "foo",
    "bar",
  ]
}
`

func testAccNetworkingRouterDataSourceName() string {
	return fmt.Sprintf(`
%s

data "vkcs_networking_router" "router" {
  name           = "${vkcs_networking_router.router.name}"
  description    = "description"
  admin_state_up = "true"
  tags = [
    "foo",
  ]
}
`, testAccNetworkingRouterDataSourceRouter)
}
