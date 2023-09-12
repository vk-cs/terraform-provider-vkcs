package dc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDCBGPNeighbor_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDCBGPNeighborBasic, map[string]string{
					"TestAccDCBGPInstanceBasic": acctest.AccTestRenderConfig(testAccDCBGPInstanceBasic, map[string]string{
						"TestAccDCRouterBasic": acctest.AccTestRenderConfig(testAccDCRouterBasic),
					}),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_bgp_neighbor.dc_bgp_neighbor", "name", "tfacc-dc-bgp-neighbor"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_neighbor.dc_bgp_neighbor", "description", "tfacc-dc-bgp-neighbor-description"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_neighbor.dc_bgp_neighbor", "add_paths", "on"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_neighbor.dc_bgp_neighbor", "remote_asn", "1"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_neighbor.dc_bgp_neighbor", "remote_ip", "192.168.1.3"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_neighbor.dc_bgp_neighbor", "enabled", "true"),
				),
			},
			{
				ResourceName:      "vkcs_dc_bgp_neighbor.dc_bgp_neighbor",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDCBGPNeighborBasic = `
{{ .TestAccDCBGPInstanceBasic }}

resource "vkcs_dc_bgp_neighbor" "dc_bgp_neighbor" {
    name = "tfacc-dc-bgp-neighbor"
	description = "tfacc-dc-bgp-neighbor-description"
    add_paths = "on"
    dc_bgp_id = vkcs_dc_bgp_instance.dc_bgp_instance.id
    remote_asn = 1
    remote_ip = "192.168.1.3"
    enabled = true
}

`
