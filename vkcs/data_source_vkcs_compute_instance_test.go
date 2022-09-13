package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeInstanceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccComputeInstanceDataSourceBasic),
			},
			{
				Config: testAccRenderConfig(testAccComputeInstanceDataSourceSource, map[string]string{"TestAccComputeInstanceDataSourceSource": testAccRenderConfig(testAccComputeInstanceDataSourceBasic)}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceDataSourceID("data.vkcs_compute_instance.source_1"),
					resource.TestCheckResourceAttr("data.vkcs_compute_instance.source_1", "name", "instance_1"),
					resource.TestCheckResourceAttrPair("data.vkcs_compute_instance.source_1", "metadata", "vkcs_compute_instance.instance_1", "metadata"),
					resource.TestCheckResourceAttrSet("data.vkcs_compute_instance.source_1", "network.0.name"),
				),
			},
		},
	})
}

func testAccCheckComputeInstanceDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find compute instance data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Compute instance data source ID not set")
		}

		return nil
	}
}

const testAccComputeInstanceDataSourceBasic = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  security_groups = ["default"]
  metadata = {
    foo = "bar"
  }
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}
`

const testAccComputeInstanceDataSourceSource = `
{{.TestAccComputeInstanceDataSourceBasic}}

data "vkcs_compute_instance" "source_1" {
  id = vkcs_compute_instance.instance_1.id
}
`
