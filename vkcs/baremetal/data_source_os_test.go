package baremetal_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccBareMetalOSDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccOSDataSourceBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vkcs_baremetal_os.basic", "id"),
				),
			},
		},
	})
}

const testAccOSDataSourceBasic = `
data "vkcs_baremetal_os" "basic" {
  name      = "ubuntu"
  version   = "22.04"
  raid_type = "no_raid"
}
`
