package images_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccImagesImagesDataSource_basic(t *testing.T) {
	testRunID := "tfacc-basic-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)

	baseConfig := acctest.AccTestRenderConfig(
		testAccImagesDataSourceBase,
		map[string]string{
			"TestRunID": testRunID,
		},
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
					testAccImagesDataSourceBasic,
					map[string]string{
						"TestAccImagesDataSourceBase": baseConfig,
						"TestRunID":                   testRunID,
					},
				),
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
	testRunID := "tfacc-filters-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)

	createdAt := time.Now().Add(-time.Hour).Format(time.RFC3339)
	updatedAt := time.Now().Add(time.Hour).Format(time.RFC3339)

	baseConfig := acctest.AccTestRenderConfig(
		testAccImagesDataSourceBase,
		map[string]string{
			"TestRunID": testRunID,
		},
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
					testAccImagesDataSourceFilterTags,
					map[string]string{
						"TestAccImagesDataSourceBase": baseConfig,
						"TestRunID":                   testRunID,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_images_images.images", "images.#", "2"),
					testAccCheckImagesDataSourceImagesNames("data.vkcs_images_images.images",
						[]string{"Centos-tf_1", "Centos-tf_2"}),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(
					testAccImagesDataSourceFilterProperties,
					map[string]string{
						"TestAccImagesDataSourceBase": baseConfig,
						"TestRunID":                   testRunID,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_images_images.images", "images.#", "2"),
					testAccCheckImagesDataSourceImagesNames("data.vkcs_images_images.images",
						[]string{"CirrOS-tf_1", "Centos-tf_2"}),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(
					testAccImagesDataSourceFilterDate,
					map[string]string{
						"TestAccImagesDataSourceBase": baseConfig,
						"CreatedAt":                   createdAt,
						"UpdatedAt":                   updatedAt,
						"TestRunID":                   testRunID,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_images_images.images", "images.#", "3"),
				),
			},
		},
	})
}

func TestAccImagesImagesDataSource_default(t *testing.T) {
	testRunID := "tfacc-default-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)

	baseConfig := acctest.AccTestRenderConfig(
		testAccImagesDataSourceDefaultBase,
		map[string]string{
			"TestRunID": testRunID,
		},
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
					testAccImagesDataSourceFilterDefault,
					map[string]string{
						"TestAccImagesDataSourceDefaultBase": baseConfig,
						"TestRunID":                          testRunID,
					},
				),
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
  tags = [
    "cirros",
    "{{.TestRunID}}",
  ]
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
  tags = [
    "centos",
    "{{.TestRunID}}",
  ]
  properties = {
	bar = "foo"
  }
}

resource "vkcs_images_image" "image_3" {
  name   = "Centos-tf_2"
  image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
  container_format = "bare"
  disk_format = "raw"
  tags = [
    "centos",
    "{{.TestRunID}}",
  ]
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
  tags       = ["{{.TestRunID}}"]
}
`

const testAccImagesDataSourceFilterTags = `
{{.TestAccImagesDataSourceBase}}

data "vkcs_images_images" "images" {
  visibility = "private"
  tags       = [
    "centos",
	"{{.TestRunID}}",
  ]
}
`

const testAccImagesDataSourceFilterProperties = `
{{.TestAccImagesDataSourceBase}}

data "vkcs_images_images" "images" {
  visibility = "private"
  properties = {
    foo = "bar"
  }
  tags = ["{{.TestRunID}}"]
}
`

const testAccImagesDataSourceFilterDate = `
{{.TestAccImagesDataSourceBase}}

data "vkcs_images_images" "images" {
  visibility = "private"
  created_at = "gt:{{ .CreatedAt }}"
  updated_at = "lte:{{ .UpdatedAt }}"
  tags       = ["{{.TestRunID}}"]
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
  tags = ["{{.TestRunID}}"]
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
  tags = ["{{.TestRunID}}"]
}

resource "vkcs_images_image" "image_3" {
  name   = "Centos-tf_2"
  image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
  container_format = "bare"
  disk_format = "raw"
  tags = ["{{.TestRunID}}"]
}
`

const testAccImagesDataSourceFilterDefault = `
{{ .TestAccImagesDataSourceDefaultBase }}

data "vkcs_images_images" "images" {
  visibility = "private"
  default    = true
  tags       = ["{{.TestRunID}}"]
}
`
