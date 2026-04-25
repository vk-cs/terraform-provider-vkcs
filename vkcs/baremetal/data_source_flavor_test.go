package baremetal_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccBareMetalFlavorDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccFlavorDataSourceBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vkcs_baremetal_flavor.basic", "cpu_model"),
					resource.TestCheckResourceAttrSet("data.vkcs_baremetal_flavor.basic", "cpu_cores"),
					resource.TestCheckResourceAttrSet("data.vkcs_baremetal_flavor.basic", "ram_size"),
					resource.TestCheckResourceAttrSet("data.vkcs_baremetal_flavor.basic", "bond_vlan_capable"),
				),
			},
		},
	})
}

const testAccFlavorDataSourceBasic = `
data "vkcs_baremetal_flavor" "basic" {
  name = "test_flavor2"
}
`
