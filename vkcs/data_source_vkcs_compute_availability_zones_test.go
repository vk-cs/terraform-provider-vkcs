package vkcs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAvailabilityZones_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAvailabilityZonesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.vkcs_compute_availability_zones.zones", "names.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

const testAccAvailabilityZonesConfig = `
data "vkcs_compute_availability_zones" "zones" {}
`
