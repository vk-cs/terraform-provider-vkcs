package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccImagesImageDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckImage(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageDataSourceCirros,
			},
			{
				Config: testAccImagesImageDataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_images_image.image_1", "name", "CirrOS-tf_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_images_image.image_1", "container_format", "bare"),
					resource.TestCheckResourceAttr(
						"data.vkcs_images_image.image_1", "disk_format", "qcow2"),
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
		PreCheck:          func() { testAccPreCheckImage(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageDataSourceCirros,
			},
			{
				Config: testAccImagesImageDataSourceQueryTag(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_1"),
				),
			},
			{
				Config: testAccImagesImageDataSourceQuerySizeMin(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_1"),
				),
			},
			{
				Config: testAccImagesImageDataSourceQuerySizeMax(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_1"),
				),
			},
			{
				Config: testAccImagesImageDataSourceQueryHidden(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_3"),
				),
			},
			{
				Config: testAccImagesImageDataSourceProperty(),
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
  disk_format = "qcow2"
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
  disk_format = "qcow2"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  tags = ["cirros-tf_2"]
  properties = {
    foo = "bar"
  }
}

resource "vkcs_images_image" "image_3" {
  name = "CirrOS-tf_3"
  container_format = "bare"
  hidden = true
  disk_format = "qcow2"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  tags = ["cirros-tf_3"]
  properties = {
	foo = "bar"
  }
}
`

func testAccImagesImageDataSourceBasic() string {
	return fmt.Sprintf(`
%s

data "vkcs_images_image" "image_1" {
	most_recent = true
	name = "${vkcs_images_image.image_1.name}"
}
`, testAccImagesImageDataSourceCirros)
}

func testAccImagesImageDataSourceQueryTag() string {
	return fmt.Sprintf(`
%s

data "vkcs_images_image" "image_1" {
	most_recent = true
	visibility = "private"
	tag = "cirros-tf_1"
}
`, testAccImagesImageDataSourceCirros)
}

func testAccImagesImageDataSourceQuerySizeMin() string {
	return fmt.Sprintf(`
%s

data "vkcs_images_image" "image_1" {
	most_recent = true
	visibility = "private"
	size_min = "13000000"
}
`, testAccImagesImageDataSourceCirros)
}

func testAccImagesImageDataSourceQuerySizeMax() string {
	return fmt.Sprintf(`
%s

data "vkcs_images_image" "image_1" {
	most_recent = true
	visibility = "private"
	size_max = "23000000"
}
`, testAccImagesImageDataSourceCirros)
}

func testAccImagesImageDataSourceQueryHidden() string {
	return fmt.Sprintf(`
%s

data "vkcs_images_image" "image_3" {
	most_recent = true
	visibility = "private"
	hidden = true
}
`, testAccImagesImageDataSourceCirros)
}

func testAccImagesImageDataSourceProperty() string {
	return fmt.Sprintf(`
%s

data "vkcs_images_image" "image_1" {
  properties = {
    foo = "bar"
    bar = "foo"
  }
}
`, testAccImagesImageDataSourceCirros)
}
