package dc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDCAPIOptions_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDCAPIOptionsBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vkcs_dc_api_options.dc_api_options", "availability_zones.#"),
					resource.TestCheckResourceAttrSet("data.vkcs_dc_api_options.dc_api_options", "flavors.#"),
				),
			},
		},
	})
}

const testAccDCAPIOptionsBasic = `
data "vkcs_dc_api_options" "dc_api_options" {}
`
