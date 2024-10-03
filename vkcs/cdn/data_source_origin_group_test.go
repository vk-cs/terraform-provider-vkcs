package cdn_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccCDNOriginGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCDNOriginGroupDataSourceBasic, map[string]string{"TestAccCDNOriginGroupDataSourceBase": testAccCDNOriginGroupDataSourceBase}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_cdn_origin_group.basic", "name", "tfacc-origin-group-base"),
					resource.TestCheckResourceAttrPair("data.vkcs_cdn_origin_group.basic", "origins", "vkcs_cdn_origin_group.base", "origins"),
					resource.TestCheckResourceAttrPair("data.vkcs_cdn_origin_group.basic", "use_next", "vkcs_cdn_origin_group.base", "use_next"),
				),
			},
		},
	})
}

const testAccCDNOriginGroupDataSourceBase = `
resource "vkcs_cdn_origin_group" "base" {
  name = "tfacc-origin-group-base"
  origins = [
    {
      source = "example.com"
    }
  ]
}
`

const testAccCDNOriginGroupDataSourceBasic = `
{{ .TestAccCDNOriginGroupDataSourceBase }}

data "vkcs_cdn_origin_group" "basic" {
  name = vkcs_cdn_origin_group.base.name
}
`
