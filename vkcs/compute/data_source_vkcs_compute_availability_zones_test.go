package compute_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccComputeAvailabilityZones_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
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
