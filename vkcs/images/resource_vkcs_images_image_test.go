package images_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

func TestAccImagesImage_basic(t *testing.T) {
	var image images.Image

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckImagesImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "name", "Centos TerraformAccTest"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "container_format", "bare"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "disk_format", "raw"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "schema", "/v2/schemas/image"),
				),
			},
		},
	})
}

func TestAccImagesImage_name(t *testing.T) {
	var image images.Image

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckImagesImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageName1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "name", "Centos TerraformAccTest"),
				),
			},
			{
				Config: testAccImagesImageName2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "name", "TerraformAccTest Centos"),
				),
			},
		},
	})
}

func TestAccImagesImage_tags(t *testing.T) {
	var image images.Image

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckImagesImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageTags1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image),
					testAccCheckImagesImageHasTag("vkcs_images_image.image_1", "foo"),
					testAccCheckImagesImageHasTag("vkcs_images_image.image_1", "bar"),
					testAccCheckImagesImageTagCount("vkcs_images_image.image_1", 2),
				),
			},
			{
				Config: testAccImagesImageTags2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image),
					testAccCheckImagesImageHasTag("vkcs_images_image.image_1", "foo"),
					testAccCheckImagesImageHasTag("vkcs_images_image.image_1", "bar"),
					testAccCheckImagesImageHasTag("vkcs_images_image.image_1", "baz"),
					testAccCheckImagesImageTagCount("vkcs_images_image.image_1", 3),
				),
			},
			{
				Config: testAccImagesImageTags3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image),
					testAccCheckImagesImageHasTag("vkcs_images_image.image_1", "foo"),
					testAccCheckImagesImageHasTag("vkcs_images_image.image_1", "baz"),
					testAccCheckImagesImageTagCount("vkcs_images_image.image_1", 2),
				),
			},
		},
	})
}

func TestAccImagesImage_properties(t *testing.T) {
	var image1 images.Image
	var image2 images.Image
	var image3 images.Image
	var image4 images.Image
	var image5 images.Image

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckImagesImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image1),
				),
			},
			{
				Config: testAccImagesImageProperties1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image2),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "properties.foo", "bar"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "properties.bar", "foo"),
				),
			},
			{
				Config: testAccImagesImageProperties2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image3),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "properties.foo", "bar"),
				),
			},
			{
				Config: testAccImagesImageProperties3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image4),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "properties.foo", "baz"),
				),
			},
			{
				Config: testAccImagesImageProperties4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image5),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "properties.foo", "baz"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "properties.bar", "foo"),
				),
			},
		},
	})
}

func testAccCheckImagesImageDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	imageClient, err := config.ImageV2Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS Image: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_images_image" {
			continue
		}

		_, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Image still exists")
		}
	}

	return nil
}

func testAccCheckImagesImageExists(n string, image *images.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		imageClient, err := config.ImageV2Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS Image: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Image not found")
		}

		*image = *found

		return nil
	}
}

func testAccCheckImagesImageHasTag(n, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		imageClient, err := config.ImageV2Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS Image: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Image not found")
		}

		for _, v := range found.Tags {
			if tag == v {
				return nil
			}
		}

		return fmt.Errorf("Tag not found: %s", tag)
	}
}

func testAccCheckImagesImageTagCount(n string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		imageClient, err := config.ImageV2Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS Image: %s", err)
		}

		found, err := images.Get(imageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Image not found")
		}

		if len(found.Tags) != expected {
			return fmt.Errorf("Expecting %d tags, found %d", expected, len(found.Tags))
		}

		return nil
	}
}

const testAccImagesImageBasic = `
  resource "vkcs_images_image" "image_1" {
      name   = "Centos TerraformAccTest"
      image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
      container_format = "bare"
      disk_format = "raw"

      timeouts {
        create = "10m"
      }
  }`

const testAccImagesImageName1 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Centos TerraformAccTest"
      image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
      container_format = "bare"
      disk_format = "raw"
  }`

const testAccImagesImageName2 = `
  resource "vkcs_images_image" "image_1" {
      name   = "TerraformAccTest Centos"
      image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
      container_format = "bare"
      disk_format = "raw"
  }`

const testAccImagesImageTags1 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Centos TerraformAccTest"
      image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
      container_format = "bare"
      disk_format = "raw"
      tags = ["foo","bar"]
  }`

const testAccImagesImageTags2 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Centos TerraformAccTest"
      image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
      container_format = "bare"
      disk_format = "raw"
      tags = ["foo","bar","baz"]
  }`

const testAccImagesImageTags3 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Centos TerraformAccTest"
      image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
      container_format = "bare"
      disk_format = "raw"
      tags = ["foo","baz"]
  }`

const testAccImagesImageProperties1 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Centos TerraformAccTest"
      image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
      container_format = "bare"
      disk_format = "raw"

      properties = {
        foo = "bar"
        bar = "foo"
      }
  }`

const testAccImagesImageProperties2 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Centos TerraformAccTest"
      image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
      container_format = "bare"
      disk_format = "raw"

      properties = {
        foo = "bar"
      }
  }`

const testAccImagesImageProperties3 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Centos TerraformAccTest"
      image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
      container_format = "bare"
      disk_format = "raw"

      properties = {
        foo = "baz"
      }
  }`

const testAccImagesImageProperties4 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Centos TerraformAccTest"
      image_source_url = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.raw.tar.gz"
      container_format = "bare"
      disk_format = "raw"

      properties = {
        foo = "baz"
        bar = "foo"
      }
  }`
