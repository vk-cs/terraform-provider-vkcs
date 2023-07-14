package images_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccImagesImagesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccImagesDataSourceBase),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccImagesDataSourceBasic, map[string]string{"TestAccImagesDataSourceBase": testAccImagesDataSourceBase}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_images_images.images", "images.#", "3"),
					resource.TestCheckResourceAttrSet("data.vkcs_images_images.images", "images.0.id"),
					resource.TestCheckResourceAttrSet("data.vkcs_images_images.images", "images.0.name"),
				),
			},
		},
	})
}

func TestAccImagesImagesDataSource_filters(t *testing.T) {
	createdAt := time.Now().Add(-time.Hour).Format(time.RFC3339)
	updatedAt := time.Now().Add(time.Hour).Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccImagesDataSourceBase),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccImagesDataSourceFilterTags, map[string]string{"TestAccImagesDataSourceBase": testAccImagesDataSourceBase}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_images_images.images", "images.#", "2"),
					testAccCheckImagesDataSourceImagesNames("data.vkcs_images_images.images",
						[]string{"Centos-tf_1", "Centos-tf_2"}),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccImagesDataSourceFilterProperties, map[string]string{"TestAccImagesDataSourceBase": testAccImagesDataSourceBase}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_images_images.images", "images.#", "2"),
					testAccCheckImagesDataSourceImagesNames("data.vkcs_images_images.images",
						[]string{"CirrOS-tf_1", "Centos-tf_2"}),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccImagesDataSourceFilterDate, map[string]string{"TestAccImagesDataSourceBase": testAccImagesDataSourceBase, "CreatedAt": createdAt, "UpdatedAt": updatedAt}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_images_images.images", "images.#", "3"),
				),
			},
		},
	})
}

func TestAccImagesImagesDataSource_default(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccImagesDataSourceDefaultBase),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccImagesDataSourceFilterDefault, map[string]string{"TestAccImagesDataSourceDefaultBase": testAccImagesDataSourceDefaultBase}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_images_images.images", "images.#", "2"),
					testAccCheckImagesDataSourceImagesNames("data.vkcs_images_images.images",
						[]string{"Centos-tf_1", "Centos-tf_2"}),
				),
			},
		},
	})
}

func testAccCheckImagesDataSourceImagesNames(n string, expectedNames []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find image data source: %s", n)
		}

		expNamesMap := make(map[string]bool)
		for _, n := range expectedNames {
			expNamesMap[n] = true
		}

		n, _ := strconv.Atoi(rs.Primary.Attributes["images.#"])
		for i := 0; i < n; i++ {
			name := rs.Primary.Attributes[fmt.Sprintf("images.%d.name", i)]
			if _, ok := expNamesMap[name]; !ok {
				return fmt.Errorf("image name %s is not in expected values: %v", name, expectedNames)
			}
		}

		return nil
	}
}

const testAccImagesDataSourceBase = `
resource "vkcs_images_image" "image_1" {
  name = "CirrOS-tf_1"
  container_format = "bare"
  disk_format = "raw"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  tags = ["cirros"]
  properties = {
	foo = "bar"
	bar = "foo"
  }
}

resource "vkcs_images_image" "image_2" {
  name   = "Centos-tf_1"
  image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
  container_format = "bare"
  disk_format = "raw"
  tags = ["centos"]
  properties = {
	bar = "foo"
  }
}

resource "vkcs_images_image" "image_3" {
  name   = "Centos-tf_2"
  image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
  container_format = "bare"
  disk_format = "raw"
  tags = ["centos"]
  properties = {
	foo = "bar"
	bar = "foo"
	foobar = "barfoo"
  }
}
`

const testAccImagesDataSourceBasic = `
{{.TestAccImagesDataSourceBase}}

data "vkcs_images_images" "images" {
  visibility = "private"
}
`

const testAccImagesDataSourceFilterTags = `
{{.TestAccImagesDataSourceBase}}

data "vkcs_images_images" "images" {
  visibility = "private"
  tags       = ["centos"]
}
`

const testAccImagesDataSourceFilterProperties = `
{{.TestAccImagesDataSourceBase}}

data "vkcs_images_images" "images" {
  visibility = "private"
  properties = {
    foo = "bar"
  }
}
`

const testAccImagesDataSourceFilterDate = `
{{.TestAccImagesDataSourceBase}}

data "vkcs_images_images" "images" {
  visibility = "private"
  created_at = "gt:{{ .CreatedAt }}"
  updated_at = "lte:{{ .UpdatedAt }}"
}
`

const testAccImagesDataSourceDefaultBase = `
resource "vkcs_images_image" "image_1" {
  name = "CirrOS-tf_1"
  container_format = "bare"
  disk_format = "raw"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  properties = {
	image_type = "snapshot"
  }
}

resource "vkcs_images_image" "image_2" {
  name   = "Centos-tf_1"
  image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
  container_format = "bare"
  disk_format = "raw"
  properties = {
	sid = "ml"
	image_type = "image"
  }
}

resource "vkcs_images_image" "image_3" {
  name   = "Centos-tf_2"
  image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
  container_format = "bare"
  disk_format = "raw"
}
`

const testAccImagesDataSourceFilterDefault = `
{{ .TestAccImagesDataSourceDefaultBase }}

data "vkcs_images_images" "images" {
  visibility = "private"
  default    = true
}
`
