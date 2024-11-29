package monitoring_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccCloudMonitoring_basic(t *testing.T) {
	resourceName := "vkcs_cloud_monitoring.basic"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCloudMonitoringBase, map[string]string{"BaseImageDataSource": testAccImageUbuntuDataSource}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "service_user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "script"),
				),
			},
		},
	})
}

const testAccCloudMonitoringBase = `
{{.BaseImageDataSource}}

resource "vkcs_cloud_monitoring" "basic" {
  image_id = data.vkcs_images_image.basic.id
}
`

const testAccImageUbuntuDataSource = `
data "vkcs_images_image" "basic" {
  visibility = "public"
  default    = true
  properties = {
    mcs_os_distro  = "ubuntu"
    mcs_os_version = "24.04"
  }
}
`
