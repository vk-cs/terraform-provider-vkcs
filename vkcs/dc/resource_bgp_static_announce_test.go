package dc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDCBGPStaticAnnounce_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDCBGPStaticAnnounceBasic, map[string]string{
					"TestAccDCBGPInstanceBasic": acctest.AccTestRenderConfig(testAccDCBGPInstanceBasic, map[string]string{
						"TestAccDCRouterBasic": acctest.AccTestRenderConfig(testAccDCRouterBasic),
					}),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_bgp_static_announce.dc_bgp_static_announce", "name", "tfacc-dc-bgp-static-announce"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_static_announce.dc_bgp_static_announce", "description", "tfacc-dc-bgp-static-announce-description"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_static_announce.dc_bgp_static_announce", "network", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_static_announce.dc_bgp_static_announce", "gateway", "192.168.1.3"),
					resource.TestCheckResourceAttr("vkcs_dc_bgp_static_announce.dc_bgp_static_announce", "enabled", "true"),
				),
			},
			{
				ResourceName:      "vkcs_dc_bgp_static_announce.dc_bgp_static_announce",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDCBGPStaticAnnounceBasic = `
{{ .TestAccDCBGPInstanceBasic }}

resource "vkcs_dc_bgp_static_announce" "dc_bgp_static_announce" {
    name = "tfacc-dc-bgp-static-announce"
    description = "tfacc-dc-bgp-static-announce-description"
    dc_bgp_id = vkcs_dc_bgp_instance.dc_bgp_instance.id
    network = "192.168.1.0/24"
    gateway = "192.168.1.3"
    enabled = true
}

`
