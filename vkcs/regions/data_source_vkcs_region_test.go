package regions_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccRegionDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_region.default", "id", "RegionOne"),
					resource.TestCheckResourceAttr("data.vkcs_region.default", "description", ""),
					resource.TestCheckResourceAttr("data.vkcs_region.default", "parent_region", ""),
				),
			},
		},
	})
}

func TestAccRegionDataSource_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionDataSourceConfigID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_region.id", "id", "RegionAms"),
					resource.TestCheckResourceAttr("data.vkcs_region.id", "description", ""),
					resource.TestCheckResourceAttr("data.vkcs_region.id", "parent_region", ""),
				),
			},
		},
	})
}

const testAccRegionDataSourceConfigBasic = `
data "vkcs_region" "default" {}
`

const testAccRegionDataSourceConfigID = `
data "vkcs_region" "id" {
	id = "RegionAms"
}
`
