package firewall_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccFirewallSecGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallSecGroupDataSourceGroup,
			},
			{
				Config: acctest.AccTestRenderConfig(testAccFirewallSecGroupDataSourceBasic, map[string]string{"TestAccFirewallSecGroupDataSourceGroup": testAccFirewallSecGroupDataSourceGroup}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupDataSourceID("data.vkcs_networking_secgroup.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_secgroup.secgroup_1", "name", "secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_secgroup.secgroup_1", "description", "My neutron security group"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_secgroup.secgroup_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_secgroup.secgroup_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

func TestAccFirewallSecGroupDataSource_secGroupID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallSecGroupDataSourceGroup,
			},
			{
				Config: acctest.AccTestRenderConfig(testAccFirewallSecGroupDataSourceSecGroupID, map[string]string{"TestAccFirewallSecGroupDataSourceGroup": testAccFirewallSecGroupDataSourceGroup}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupDataSourceID("data.vkcs_networking_secgroup.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_secgroup.secgroup_1", "name", "secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_secgroup.secgroup_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_secgroup.secgroup_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSecGroupDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find security group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Security group data source ID not set")
		}

		return nil
	}
}

const testAccFirewallSecGroupDataSourceGroup = `
resource "vkcs_networking_secgroup" "secgroup_1" {
  name        = "secgroup_1"
  description = "My neutron security group"
  tags = [
    "foo",
    "bar",
  ]
}
`

const testAccFirewallSecGroupDataSourceBasic = `
{{.TestAccFirewallSecGroupDataSourceGroup}}

data "vkcs_networking_secgroup" "secgroup_1" {
  name = vkcs_networking_secgroup.secgroup_1.name
  tags = [
    "bar",
  ]
}
`

const testAccFirewallSecGroupDataSourceSecGroupID = `
{{.TestAccFirewallSecGroupDataSourceGroup}}

data "vkcs_networking_secgroup" "secgroup_1" {
  secgroup_id = vkcs_networking_secgroup.secgroup_1.id
  tags = [
    "foo",
  ]
}
`
