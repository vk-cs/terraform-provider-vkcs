package images_test

import (
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccImagesImageDataSource_basic(t *testing.T) {
	testRunID := "tfacc-image-basic-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)

	baseConfig := acctest.AccTestRenderConfig(
		testAccImagesImageDataSourceCirros,
		map[string]string{"TestRunID": testRunID},
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: baseConfig,
			},
			{
				Config: acctest.AccTestRenderConfig(
					testAccImagesImageDataSourceBasic,
					map[string]string{
						"TestAccImagesImageDataSourceCirros": baseConfig,
						"TestRunID":                          testRunID,
					},
				),
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
	testRunID := "tfacc-image-queries-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)

	baseConfig := acctest.AccTestRenderConfig(
		testAccImagesImageDataSourceCirros,
		map[string]string{"TestRunID": testRunID},
	)

	renderQuery := func(query string) string {
		return acctest.AccTestRenderConfig(
			query,
			map[string]string{
				"TestAccImagesImageDataSourceCirros": baseConfig,
				"TestRunID":                          testRunID,
			},
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: baseConfig,
			},
			{
				Config: renderQuery(testAccImagesImageDataSourceQueryTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_1"),
				),
			},
			{
				Config: renderQuery(testAccImagesImageDataSourceQuerySizeMin),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_1"),
				),
			},
			{
				Config: renderQuery(testAccImagesImageDataSourceQuerySizeMax),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_1"),
				),
			},
			{
				Config: renderQuery(testAccImagesImageDataSourceProperty),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_1"),
				),
			},
			{
				Config: baseConfig,
			},
		},
	})
}

func TestAccImagesImageDataSource_default(t *testing.T) {
	testRunID := "tfacc-image-default-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)

	baseConfig := acctest.AccTestRenderConfig(
		testAccImagesImageDataSourceDefaultBase,
		map[string]string{"TestRunID": testRunID},
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: baseConfig,
			},
			{
				Config: acctest.AccTestRenderConfig(
					testAccImagesImageDataSourceDefault,
					map[string]string{
						"TestAccImagesImageDataSourceDefaultBase": baseConfig,
						"TestRunID": testRunID,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesDataSourceID("data.vkcs_images_image.image_1"),
					resource.TestCheckResourceAttr(
						"data.vkcs_images_image.image_1", "name", "CirrOS-tf_3"),
				),
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
  tags = ["{{.TestRunID}}-1"]
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
  tags = ["{{.TestRunID}}-2"]
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
	tag = "{{.TestRunID}}-1"
}
`

const testAccImagesImageDataSourceQueryTag = `
{{.TestAccImagesImageDataSourceCirros}}

data "vkcs_images_image" "image_1" {
	most_recent = true
	visibility = "private"
	tag = "{{.TestRunID}}-1"
}
`

const testAccImagesImageDataSourceQuerySizeMin = `
{{.TestAccImagesImageDataSourceCirros}}

data "vkcs_images_image" "image_1" {
	most_recent = true
	visibility = "private"
	size_min = "13000000"
	tag = "{{.TestRunID}}-1"
}
`

const testAccImagesImageDataSourceQuerySizeMax = `
{{.TestAccImagesImageDataSourceCirros}}

data "vkcs_images_image" "image_1" {
	most_recent = true
	visibility = "private"
	size_max = "23000000"
	tag = "{{.TestRunID}}-1"
}
`

const testAccImagesImageDataSourceProperty = `
{{.TestAccImagesImageDataSourceCirros}}

data "vkcs_images_image" "image_1" {
  tag = "{{.TestRunID}}-1"
  properties = {
    foo = "bar"
    bar = "foo"
  }
}
`

const testAccImagesImageDataSourceDefaultBase = `
resource "vkcs_images_image" "image_1" {
  name = "CirrOS-tf_1"
  container_format = "bare"
  disk_format = "raw"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  tags = ["{{.TestRunID}}"]
  properties = {
	sid = "hidden"
  }
}

resource "vkcs_images_image" "image_2" {
  name = "CirrOS-tf_2"
  container_format = "bare"
  disk_format = "raw"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  tags = ["{{.TestRunID}}"]
  properties = {
    image_type = "snapshot"
  }
}

resource "vkcs_images_image" "image_3" {
  name = "CirrOS-tf_3"
  container_format = "bare"
  disk_format = "raw"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  tags = ["{{.TestRunID}}"]
}
`

const testAccImagesImageDataSourceDefault = `
{{.TestAccImagesImageDataSourceDefaultBase}}

data "vkcs_images_image" "image_1" {
  visibility = "private"
  default = true
  tag = "{{.TestRunID}}"
}
`
