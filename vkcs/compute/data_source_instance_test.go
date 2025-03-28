package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccComputeInstanceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceDataSourceBasic),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeInstanceDataSourceSource, map[string]string{"TestAccComputeInstanceDataSourceBasic": acctest.AccTestRenderConfig(testAccComputeInstanceDataSourceBasic)}),
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
{{.BaseSecurityGroup}}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
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
