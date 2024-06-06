package dc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDCConntrackHelper_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDCConntrackHelperBasic, map[string]string{
					"TestAccDCRouterBasic": acctest.AccTestRenderConfig(testAccDCRouterBasic),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_conntrack_helper.dc_conntrack_helper", "helper", "ftp"),
					resource.TestCheckResourceAttr("vkcs_dc_conntrack_helper.dc_conntrack_helper", "protocol", "tcp"),
					resource.TestCheckResourceAttr("vkcs_dc_conntrack_helper.dc_conntrack_helper", "port", "80"),
				),
			},
			{
				ResourceName:      "vkcs_dc_conntrack_helper.dc_conntrack_helper",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDCConntrackHelper_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDCConntrackHelperFull, map[string]string{
					"TestAccDCRouterBasic": acctest.AccTestRenderConfig(testAccDCRouterBasic),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_conntrack_helper.dc_conntrack_helper", "name", "tfacc-dc-conntrack-helper"),
					resource.TestCheckResourceAttr("vkcs_dc_conntrack_helper.dc_conntrack_helper", "description", "tfacc-dc-conntrack-helper-description"),
					resource.TestCheckResourceAttr("vkcs_dc_conntrack_helper.dc_conntrack_helper", "helper", "ftp"),
					resource.TestCheckResourceAttr("vkcs_dc_conntrack_helper.dc_conntrack_helper", "protocol", "tcp"),
					resource.TestCheckResourceAttr("vkcs_dc_conntrack_helper.dc_conntrack_helper", "port", "80"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDCConntrackHelperFullUpdate, map[string]string{
					"TestAccDCRouterBasic": acctest.AccTestRenderConfig(testAccDCRouterBasic),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_conntrack_helper.dc_conntrack_helper", "name", "tfacc-dc-conntrack-helper-upd"),
					resource.TestCheckResourceAttr("vkcs_dc_conntrack_helper.dc_conntrack_helper", "description", "tfacc-dc-conntrack-helper-description-upd"),
					resource.TestCheckResourceAttr("vkcs_dc_conntrack_helper.dc_conntrack_helper", "helper", "ftp"),
					resource.TestCheckResourceAttr("vkcs_dc_conntrack_helper.dc_conntrack_helper", "protocol", "tcp"),
					resource.TestCheckResourceAttr("vkcs_dc_conntrack_helper.dc_conntrack_helper", "port", "81"),
				),
			},
		},
	})
}

const testAccDCConntrackHelperBasic = `
{{ .TestAccDCRouterBasic }}

resource "vkcs_dc_conntrack_helper" "dc_conntrack_helper" {
	dc_router_id = vkcs_dc_router.dc_router.id
	helper = "ftp"
    protocol = "tcp"
    port = 80
}
`

const testAccDCConntrackHelperFull = `
{{ .TestAccDCRouterBasic }}

resource "vkcs_dc_conntrack_helper" "dc_conntrack_helper" {
	name = "tfacc-dc-conntrack-helper"
	description = "tfacc-dc-conntrack-helper-description"
	dc_router_id = vkcs_dc_router.dc_router.id
	helper = "ftp"
    protocol = "tcp"
    port = 80
}
`

const testAccDCConntrackHelperFullUpdate = `
{{ .TestAccDCRouterBasic }}

resource "vkcs_dc_conntrack_helper" "dc_conntrack_helper" {
	name = "tfacc-dc-conntrack-helper-upd"
	description = "tfacc-dc-conntrack-helper-description-upd"
	dc_router_id = vkcs_dc_router.dc_router.id
	helper = "ftp"
    protocol = "tcp"
    port = 81
}
`
