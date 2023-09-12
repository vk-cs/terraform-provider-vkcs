package dc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDCRouter_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDCRouterBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_router.dc_router", "name", "tfacc-dc-router"),
					resource.TestCheckResourceAttr("vkcs_dc_router.dc_router", "description", "tfacc-dc-router-description"),
				),
			},
			{
				ResourceName:      "vkcs_dc_router.dc_router",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDCRouterBasic = `
resource "vkcs_dc_router" "dc_router" {
	name = "tfacc-dc-router"
	description = "tfacc-dc-router-description"
	availability_zone = "{{.AvailabilityZone}}"
	flavor = "standard"
  }
`
