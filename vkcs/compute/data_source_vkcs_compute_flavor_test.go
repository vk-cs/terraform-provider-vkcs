package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccComputeFlavorDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFlavorDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeFlavorDataSourceID("data.vkcs_compute_flavor.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "name", "Basic-1-2-20"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "ram", "2048"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "disk", "20"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "vcpus", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "rx_tx_factor", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "is_public", "true"),
				),
			},
		},
	})
}

func TestAccComputeFlavorDataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFlavorDataSourceQueryDisk,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeFlavorDataSourceID("data.vkcs_compute_flavor.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "name", "Basic-1-1-10"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "ram", "1024"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "disk", "10"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "vcpus", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "rx_tx_factor", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "is_public", "true"),
				),
			},
			{
				Config: testAccComputeFlavorDataSourceQueryMinDisk,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeFlavorDataSourceID("data.vkcs_compute_flavor.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "name", "Basic-1-2-20"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "ram", "2048"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "disk", "20"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "vcpus", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "rx_tx_factor", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "is_public", "true"),
				),
			},
			{
				Config: testAccComputeFlavorDataSourceQueryMinRAM,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeFlavorDataSourceID("data.vkcs_compute_flavor.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "name", "Basic-1-2-20"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "ram", "2048"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "disk", "20"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "vcpus", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "rx_tx_factor", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "is_public", "true"),
				),
			},
			{
				Config: testAccComputeFlavorDataSourceQueryVCPUs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeFlavorDataSourceID("data.vkcs_compute_flavor.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "name", "Basic-1-2-20"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "ram", "2048"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "disk", "20"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "vcpus", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "rx_tx_factor", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "is_public", "true"),
				),
			},
		},
	})
}

func TestAccComputeFlavorDataSource_extraSpecs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFlavorDataSourceExtraSpecs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeFlavorDataSourceID("data.vkcs_compute_flavor.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "name", "Basic-1-2-20"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "extra_specs.%", "3"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "extra_specs.agg_common", "true"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "extra_specs.hw:cpu_sockets", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "extra_specs.mcs:cpu_type", "standard"),
				),
			},
		},
	})
}

func testAccCheckComputeFlavorDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find flavor data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Flavor data source ID not set")
		}

		return nil
	}
}

const testAccComputeFlavorDataSourceBasic = `
data "vkcs_compute_flavor" "flavor_1" {
  name = "Basic-1-2-20"
}
`

const testAccComputeFlavorDataSourceQueryDisk = `
data "vkcs_compute_flavor" "flavor_1" {
  disk = 10
}
`

const testAccComputeFlavorDataSourceQueryMinDisk = `
data "vkcs_compute_flavor" "flavor_1" {
  name = "Basic-1-2-20"
  min_disk = 20
}
`

const testAccComputeFlavorDataSourceQueryMinRAM = `
data "vkcs_compute_flavor" "flavor_1" {
  name = "Basic-1-2-20"
  min_ram = 2048
}
`

const testAccComputeFlavorDataSourceQueryVCPUs = `
data "vkcs_compute_flavor" "flavor_1" {
  name = "Basic-1-2-20"
  vcpus = 1
}
`

const testAccComputeFlavorDataSourceExtraSpecs = `
data "vkcs_compute_flavor" "flavor_1" {
	name = "Basic-1-2-20"
  }
`
