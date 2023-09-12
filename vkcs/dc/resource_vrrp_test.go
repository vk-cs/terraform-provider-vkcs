package dc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDCVRRP_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDCVRRPBasic, map[string]string{
					"TestAccDCVRRPBase": testAccDCVRRPBase,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_vrrp.dc_vrrp", "name", "tfacc-dc-vrrp"),
					resource.TestCheckResourceAttr("vkcs_dc_vrrp.dc_vrrp", "description", "tfacc-dc-vrrp-description"),
					resource.TestCheckResourceAttr("vkcs_dc_vrrp.dc_vrrp", "group_id", "100"),
					resource.TestCheckResourceAttr("vkcs_dc_vrrp.dc_vrrp", "advert_interval", "1"),
					resource.TestCheckResourceAttr("vkcs_dc_vrrp.dc_vrrp", "enabled", "true"),
				),
			},
			{
				ResourceName:      "vkcs_dc_vrrp.dc_vrrp",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDCVRRPBase = `
resource "vkcs_dc_vrrp" "dc_vrrp" {
    name = "tfacc-dc-vrrp"
    description = "tfacc-dc-vrrp-description"
    group_id = 100
    network_id = vkcs_networking_network.base.id
    subnet_id = vkcs_networking_subnet.base.id
    advert_interval = 1
    enabled = true
}
`

const testAccDCVRRPBasic = `
{{.BaseNetwork}}
{{ .TestAccDCVRRPBase }}
`
