package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccNetworkingFloatingIPDataSource_address(t *testing.T) {
	baseConfig := acctest.AccTestRenderConfig(testAccNetworkingFloatingIPDataSourceBase)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: baseConfig,
			},
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingFloatingIPDataSourceAddress, map[string]string{"TestAccNetworkingFloatingIPDataSourceFloatingIPBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFloatingIPDataSourceID("data.vkcs_networking_floatingip.fip_1"),
					resource.TestCheckResourceAttrSet(
						"data.vkcs_networking_floatingip.fip_1", "address"),
					resource.TestCheckResourceAttrSet(
						"data.vkcs_networking_floatingip.fip_1", "pool"),
					resource.TestCheckResourceAttrSet(
						"data.vkcs_networking_floatingip.fip_1", "status"),
					resource.TestCheckResourceAttrSet(
						"data.vkcs_networking_floatingip.fip_1", "description"),
				),
			},
		},
	})
}

func testAccCheckNetworkingFloatingIPDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find floating IP data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Floating IP data source ID not set")
		}

		return nil
	}
}

const testAccNetworkingFloatingIPDataSourceBase = `
resource "vkcs_networking_floatingip" "fip_1" {
  pool        = "{{ .ExtNetName }}"
  description = "tfacc-fip"
}
`

const testAccNetworkingFloatingIPDataSourceAddress = `
{{ .TestAccNetworkingFloatingIPDataSourceFloatingIPBase }}

data "vkcs_networking_floatingip" "fip_1" {
  address = vkcs_networking_floatingip.fip_1.address
}
`
