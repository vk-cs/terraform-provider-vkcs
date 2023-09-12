package dc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDCVRRPAddress_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDCVRRPAddressBasic, map[string]string{
					"TestAccDCVRRPBasic": acctest.AccTestRenderConfig(testAccDCVRRPBasic, map[string]string{
						"TestAccDCVRRPBase": testAccDCVRRPBase,
					}),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_vrrp_address.dc_vrrp_address", "name", "tfacc-dc-vrrp-address"),
					resource.TestCheckResourceAttr("vkcs_dc_vrrp_address.dc_vrrp_address", "description", "tfacc-dc-vrrp-address-description"),
					resource.TestCheckResourceAttr("vkcs_dc_vrrp_address.dc_vrrp_address", "ip_address", "192.168.199.42"),
				),
			},
			{
				ResourceName:      "vkcs_dc_vrrp_address.dc_vrrp_address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDCVRRPAddressBasic = `
{{ .TestAccDCVRRPBasic}}

resource "vkcs_dc_vrrp_address" "dc_vrrp_address" {
    name = "tfacc-dc-vrrp-address"
    description = "tfacc-dc-vrrp-address-description"
    dc_vrrp_id = vkcs_dc_vrrp.dc_vrrp.id
    ip_address = "192.168.199.42"
}

`
