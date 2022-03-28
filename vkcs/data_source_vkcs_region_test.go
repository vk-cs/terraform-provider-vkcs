package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKubernetesdataSourceVkcsRegion(t *testing.T) {
	tests := map[string]struct {
		name     string
		testCase resource.TestCase
	}{
		"no params": {
			name: "data.vkcs_region.empty",
			testCase: resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: `data "vkcs_region" "empty" {}`,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("data.vkcs_region.empty", "id", "RegionOne"),
							resource.TestCheckResourceAttr("data.vkcs_region.empty", "description", ""),
							resource.TestCheckResourceAttr("data.vkcs_region.empty", "parent_region", ""),
						),
					},
				},
			},
		},
		"id provided": {
			name: "data.vkcs_region.id",
			testCase: resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: `data "vkcs_region" "id" {
									id="RegionAms"
								}`,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("data.vkcs_region.id", "id", "RegionAms"),
							resource.TestCheckResourceAttr("data.vkcs_region.id", "description", ""),
							resource.TestCheckResourceAttr("data.vkcs_region.id", "parent_region", ""),
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
