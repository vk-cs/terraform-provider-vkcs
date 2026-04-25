package baremetal_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

// TestManualAccBareMetalServerDataSource_basic is a manual acceptance test.
// It cannot be executed in an automated environment because provisioning
// a vkcs_baremetal_server requires allocation of a real physical server,
// which is costly, and not suitable for routine test runs. This is too inefficient.
func TestManualAccBareMetalServerDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccServerDataSourceBasic, map[string]string{
					"TestAccBareMetalServerResourceBasic": acctest.AccTestRenderConfig(testAccServerResourceBasic),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_baremetal_server.basic", "availability_zone", "ME1"),
					resource.TestCheckResourceAttr("data.vkcs_baremetal_server.basic", "image_source", "PUBLIC"),
					resource.TestCheckResourceAttr("data.vkcs_baremetal_server.basic", "raid_type", "RAID1"),
				),
			},
		},
	})
}

const testAccServerDataSourceBasic = `
{{.TestAccBareMetalServerResourceBasic}}

data "vkcs_baremetal_server" "basic" {
   id = vkcs_baremetal_server.basic.id
}
`
