package vkcs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccMonitoringTemplateBasic = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}

resource "vkcs_compute_instance" "instance_1" {
	depends_on = ["vkcs_networking_subnet.base"]
	name = "instance_1"
	availability_zone = "{{.AvailabilityZone}}"
	security_groups = ["default"]
	metadata = {
	  foo = "bar"
	}
	network {
	  uuid = vkcs_networking_network.base.id
	}
    block_device {
       uuid                  = data.vkcs_images_image.base.id
       source_type           = "image"
       destination_type      = "volume"
       volume_type           = "high-iops"
       volume_size           = 20
       boot_index            = 0
       delete_on_termination = true
    }
	 
	flavor_id = data.vkcs_compute_flavor.base.id
}

 resource "vkcs_monitoring_template" "basic" {
  instance_id  = vkcs_compute_instance.instance_1.id
 }
`

func TestAccMonitoringTemplate_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccMonitoringTemplateBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith("vkcs_monitoring_template.basic", "script", func(value string) error {
						if !strings.Contains(value, "sudo bash") {
							return fmt.Errorf("should be script")
						}
						return nil
					}),
				),
			},
		},
	})
}
