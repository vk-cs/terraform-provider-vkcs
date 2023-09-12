package dc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDCBGPStaticRoute_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDCBGPStaticRouteBasic, map[string]string{
					"TestAccDCRouterBasic": acctest.AccTestRenderConfig(testAccDCRouterBasic),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_static_route.dc_static_route", "name", "tfacc-dc-static-route"),
					resource.TestCheckResourceAttr("vkcs_dc_static_route.dc_static_route", "description", "tfacc-dc-static-route-description"),
					resource.TestCheckResourceAttr("vkcs_dc_static_route.dc_static_route", "network", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("vkcs_dc_static_route.dc_static_route", "gateway", "192.168.1.3"),
					resource.TestCheckResourceAttr("vkcs_dc_static_route.dc_static_route", "metric", "1"),
				),
			},
			{
				ResourceName:      "vkcs_dc_static_route.dc_static_route",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDCBGPStaticRouteBasic = `
{{ .TestAccDCRouterBasic }}

resource "vkcs_dc_static_route" "dc_static_route" {
    name = "tfacc-dc-static-route"
    description = "tfacc-dc-static-route-description"
    dc_router_id = vkcs_dc_router.dc_router.id
    network = "192.168.1.0/24"
    gateway = "192.168.1.3"
    metric = 1
}

`
