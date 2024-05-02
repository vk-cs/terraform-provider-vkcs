package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
				Config: testAccComputeFlavorDataSourceQueryMinDiskWithName,
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
				Config: testAccComputeFlavorDataSourceQueryMinRAMWithName,
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
						"data.vkcs_compute_flavor.flavor_1", "name", "STD3-4-4"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "ram", "4096"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "disk", "0"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "vcpus", "4"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "rx_tx_factor", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "is_public", "true"),
				),
			},
			{
				Config: testAccComputeFlavorDataSourceQueryMinRAMAndMinDisk,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeFlavorDataSourceID("data.vkcs_compute_flavor.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "name", "Standard-2-4-40"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "ram", "4096"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "disk", "40"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "vcpus", "2"),
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
			{
				Config: testAccComputeFlavorDataSourceQueryCPUGeneration,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeFlavorDataSourceID("data.vkcs_compute_flavor.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "name", "STD2-1-2"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "ram", "2048"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "disk", "0"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "vcpus", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "rx_tx_factor", "1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "is_public", "true"),
					resource.TestCheckResourceAttr(
						"data.vkcs_compute_flavor.flavor_1", "extra_specs.mcs:cpu_generation", "cascadelake-v1"),
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
						"data.vkcs_compute_flavor.flavor_1", "extra_specs.%", "4"),
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

const testAccComputeFlavorDataSourceQueryMinDiskWithName = `
data "vkcs_compute_flavor" "flavor_1" {
  name = "Basic-1-2-20"
  min_disk = 20
}
`

const testAccComputeFlavorDataSourceQueryMinRAMWithName = `
data "vkcs_compute_flavor" "flavor_1" {
  name = "Basic-1-2-20"
  min_ram = 2048
}
`

const testAccComputeFlavorDataSourceQueryMinDisk = `
data "vkcs_compute_flavor" "flavor_1" {
  min_disk = 20
}
`

const testAccComputeFlavorDataSourceQueryMinRAM = `
data "vkcs_compute_flavor" "flavor_1" {
  min_ram = 4096
}
`

const testAccComputeFlavorDataSourceQueryMinRAMAndMinDisk = `
data "vkcs_compute_flavor" "flavor_1" {
  min_disk = 20
  min_ram = 4096
}
`

const testAccComputeFlavorDataSourceQueryVCPUs = `
data "vkcs_compute_flavor" "flavor_1" {
  name = "Basic-1-2-20"
  vcpus = 1
}
`

const testAccComputeFlavorDataSourceQueryCPUGeneration = `
data "vkcs_compute_flavor" "flavor_1" {
  name = "STD2-1-2"
  cpu_generation = "cascadelake-v1"
}
`

const testAccComputeFlavorDataSourceExtraSpecs = `
data "vkcs_compute_flavor" "flavor_1" {
	name = "Basic-1-2-20"
  }
`
