package dc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDCInterface_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDCInterfaceBasic, map[string]string{
					"TestAccDCRouterBasic": acctest.AccTestRenderConfig(testAccDCRouterBasic),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_interface.dc_interface", "name", "tfacc-dc-interface"),
					resource.TestCheckResourceAttr("vkcs_dc_interface.dc_interface", "description", "tfacc-dc-interface-description"),
					resource.TestCheckResourceAttr("vkcs_dc_interface.dc_interface", "bgp_announce_enabled", "true"),
				),
			},
			{
				ResourceName:      "vkcs_dc_interface.dc_interface",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDCInterfaceBasic = `
{{.BaseNetwork}}

{{ .TestAccDCRouterBasic }}

resource "vkcs_dc_interface" "dc_interface" {
	name = "tfacc-dc-interface"
	description = "tfacc-dc-interface-description"
	dc_router_id = vkcs_dc_router.dc_router.id
	network_id = vkcs_networking_network.base.id
	subnet_id = vkcs_networking_subnet.base.id
	bgp_announce_enabled = true
}
`
