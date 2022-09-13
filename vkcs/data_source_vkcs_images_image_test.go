package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccImagesImageDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageDataSourceCirros,
			},
			{
				Config: testAccRenderConfig(testAccImagesImageDataSourceBasic, map[string]string{"TestAccImagesImageDataSourceCirros": testAccImagesImageDataSourceCirros}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_images_image.image_1", "name", "CirrOS-tf_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_images_image.image_1", "container_format", "bare"),
					resource.TestCheckResourceAttr(
						"data.vkcs_images_image.image_1", "disk_format", "raw"),
					resource.TestCheckResourceAttr(
						"data.vkcs_images_image.image_1", "min_disk_gb", "0"),
					resource.TestCheckResourceAttr(
						"data.vkcs_images_image.image_1", "min_ram_mb", "0"),
					resource.TestCheckResourceAttr(
						"data.vkcs_images_image.image_1", "protected", "false"),
					resource.TestCheckResourceAttr(
						"data.vkcs_images_image.image_1", "visibility", "private"),
				),
			},
		},
	})
}

func TestAccImagesImageDataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageDataSourceCirros,
			},
			{
				Config: testAccRenderConfig(testAccImagesImageDataSourceQueryTag, map[string]string{"TestAccImagesImageDataSourceCirros": testAccImagesImageDataSourceCirros}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_1"),
				),
			},
			{
				Config: testAccRenderConfig(testAccImagesImageDataSourceQuerySizeMin, map[string]string{"TestAccImagesImageDataSourceCirros": testAccImagesImageDataSourceCirros}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_1"),
				),
			},
			{
				Config: testAccRenderConfig(testAccImagesImageDataSourceQuerySizeMax, map[string]string{"TestAccImagesImageDataSourceCirros": testAccImagesImageDataSourceCirros}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_1"),
				),
			},
			{
				Config: testAccRenderConfig(testAccImagesImageDataSourceProperty, map[string]string{"TestAccImagesImageDataSourceCirros": testAccImagesImageDataSourceCirros}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_1"),
				),
			},
			{
				Config: testAccImagesImageDataSourceCirros,
			},
		},
	})
}

func testAccCheckImagesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Image data source ID not set")
		}

		return nil
	}
}

// Standard CirrOS image.
const testAccImagesImageDataSourceCirros = `
resource "vkcs_images_image" "image_1" {
  name = "CirrOS-tf_1"
  container_format = "bare"
  disk_format = "raw"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  tags = ["cirros-tf_1"]
  properties = {
    foo = "bar"
    bar = "foo"
  }
}

resource "vkcs_images_image" "image_2" {
  name = "CirrOS-tf_2"
  container_format = "bare"
  disk_format = "raw"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  tags = ["cirros-tf_2"]
  properties = {
    foo = "bar"
  }
}
`

const testAccImagesImageDataSourceBasic = `
{{.TestAccImagesImageDataSourceCirros}}

data "vkcs_images_image" "image_1" {
	most_recent = true
	name = vkcs_images_image.image_1.name
}
`

const testAccImagesImageDataSourceQueryTag = `
{{.TestAccImagesImageDataSourceCirros}}

data "vkcs_images_image" "image_1" {
	most_recent = true
	visibility = "private"
	tag = "cirros-tf_1"
}
`

const testAccImagesImageDataSourceQuerySizeMin = `
{{.TestAccImagesImageDataSourceCirros}}

data "vkcs_images_image" "image_1" {
	most_recent = true
	visibility = "private"
	size_min = "13000000"
}
`

const testAccImagesImageDataSourceQuerySizeMax = `
{{.TestAccImagesImageDataSourceCirros}}

data "vkcs_images_image" "image_1" {
	most_recent = true
	visibility = "private"
	size_max = "23000000"
}
`

const testAccImagesImageDataSourceProperty = `
{{.TestAccImagesImageDataSourceCirros}}

data "vkcs_images_image" "image_1" {
  properties = {
    foo = "bar"
    bar = "foo"
  }
}
`
