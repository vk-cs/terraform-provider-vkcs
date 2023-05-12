package regions_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDataSourceVkcsRegions(t *testing.T) {
	tests := map[string]struct {
		name     string
		testCase resource.TestCase
	}{
		"no filter": {
			name: "data.vkcs_regions.empty",
			testCase: resource.TestCase{
				ProviderFactories: acctest.AccTestProviders,
				Steps: []resource.TestStep{
					{
						Config: testAccdataSourceVkcsRegionsConfigEmpty(),
						Check: resource.ComposeTestCheckFunc(
							testAccdataSourceVkcsRegionsCheck("data.vkcs_regions.empty"),
						),
					},
				},
			},
		},
		"with parent id": {
			name: "data.vkcs_regions.parent_id",
			testCase: resource.TestCase{
				ProviderFactories: acctest.AccTestProviders,
				Steps: []resource.TestStep{
					{
						Config: testAccdataSourceVkcsRegionsConfigParentID(),
						Check: resource.ComposeTestCheckFunc(
							testAccdataSourceVkcsRegionsCheck("data.vkcs_regions.parent_id"),
						),
					},
				},
			},
		},
	}

	for name := range tests {
		tt := tests[name]
		t.Run(name, func(t *testing.T) {
			resource.ParallelTest(t, tt.testCase)
		})
	}
}

func testAccdataSourceVkcsRegionsConfigEmpty() string {
	return `
data "vkcs_regions" "empty" {}
`
}

func testAccdataSourceVkcsRegionsConfigParentID() string {
	return `
data "vkcs_regions" "parent_id" {
	parent_region_id=""
}
`
}

func testAccdataSourceVkcsRegionsCheck(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("root module has no resource called %s", resourceName)
		}

		names, namesOk := rs.Primary.Attributes["names.#"]

		if !namesOk {
			return fmt.Errorf("names attribute is missing.")
		}

		namesQuantity, err := strconv.Atoi(names)

		if err != nil {
			return fmt.Errorf("error parsing names (%s) into integer: %s", names, err)
		}

		if namesQuantity == 0 {
			return fmt.Errorf("No names found, this is probably a bug.")
		}

		return nil
	}
}
