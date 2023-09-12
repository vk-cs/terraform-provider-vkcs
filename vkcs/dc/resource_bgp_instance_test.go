package dc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDCBGPInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDCBGPInstanceBasic, map[string]string{
					"TestAccDCRouterBasic": acctest.AccTestRenderConfig(testAccDCRouterBasic),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_bgp_instance.dc_bgp_instance", "name", "tfacc-dc-bgp-instance"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_instance.dc_bgp_instance", "description", "tfacc-dc-bgp-instance-description"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_instance.dc_bgp_instance", "bgp_router_id", "192.168.1.2"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_instance.dc_bgp_instance", "asn", "12345"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_instance.dc_bgp_instance", "ecmp_enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_instance.dc_bgp_instance", "enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_instance.dc_bgp_instance", "graceful_restart", "true"),
				),
			},
			{
				ResourceName:      "vkcs_dc_bgp_instance.dc_bgp_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDCBGPInstanceBasic = `
{{ .TestAccDCRouterBasic }}

resource "vkcs_dc_bgp_instance" "dc_bgp_instance" {
    name = "tfacc-dc-bgp-instance"
    description = "tfacc-dc-bgp-instance-description"
    dc_router_id = vkcs_dc_router.dc_router.id
    bgp_router_id = "192.168.1.2"
    asn = 12345
    ecmp_enabled = true
    enabled = true
    graceful_restart = true
}
`
