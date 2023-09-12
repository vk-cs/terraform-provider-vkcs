package dc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDCVRRPInterface_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDCVRRPInterfaceBasic, map[string]string{
					"TestAccDCInterfaceBasic": acctest.AccTestRenderConfig(testAccDCInterfaceBasic, map[string]string{
						"TestAccDCRouterBasic": acctest.AccTestRenderConfig(testAccDCRouterBasic),
					}),
					"TestAccDCVRRPBase": acctest.AccTestRenderConfig(testAccDCVRRPBase),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_vrrp_interface.dc_vrrp_interface", "name", "tfacc-dc-vrrp-interface"),
					resource.TestCheckResourceAttr("vkcs_dc_vrrp_interface.dc_vrrp_interface", "description", "tfacc-dc-vrrp-interface-description"),
					resource.TestCheckResourceAttr("vkcs_dc_vrrp_interface.dc_vrrp_interface", "priority", "100"),
					resource.TestCheckResourceAttr("vkcs_dc_vrrp_interface.dc_vrrp_interface", "preempt", "true"),
					resource.TestCheckResourceAttr("vkcs_dc_vrrp_interface.dc_vrrp_interface", "master", "true"),
				),
			},
			{
				ResourceName:      "vkcs_dc_vrrp_interface.dc_vrrp_interface",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDCVRRPInterfaceBasic = `
{{ .TestAccDCInterfaceBasic }}
{{ .TestAccDCVRRPBase}}

resource "vkcs_dc_vrrp_interface" "dc_vrrp_interface" {
    name = "tfacc-dc-vrrp-interface"
    description = "tfacc-dc-vrrp-interface-description"
    dc_vrrp_id = vkcs_dc_vrrp.dc_vrrp.id
    dc_interface_id = vkcs_dc_interface.dc_interface.id
    priority = 100
    preempt = true
    master = true
}
`
