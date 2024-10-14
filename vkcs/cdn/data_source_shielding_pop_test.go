package cdn_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccCDNShieldingPopDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCDNShieldingPopBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_cdn_shielding_pop.basic", "city", "Moscow-Megafon"),
					resource.TestCheckResourceAttrSet("data.vkcs_cdn_shielding_pop.basic", "country"),
					resource.TestCheckResourceAttrSet("data.vkcs_cdn_shielding_pop.basic", "datacenter"),
				),
			},
		},
	})
}

const testAccCDNShieldingPopBasic = `
data "vkcs_cdn_shielding_pop" "basic" {
  city = "Moscow-Megafon"
}
`
