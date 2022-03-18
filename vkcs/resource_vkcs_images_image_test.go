package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

func TestAccImagesImage_basic(t *testing.T) {
	var image images.Image

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckImage(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckImagesImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "name", "Rancher TerraformAccTest"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "container_format", "bare"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "disk_format", "qcow2"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "schema", "/v2/schemas/image"),
				),
			},
			{
				Config: testAccImagesImageBasicWithID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "name", "Rancher TerraformAccTest"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "container_format", "bare"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "disk_format", "qcow2"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "schema", "/v2/schemas/image"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "image_id", "c1efdf94-9a1a-4401-88b8-d616029d2551"),
				),
			},
			{
				Config: testAccImagesImageBasicHidden,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "name", "Rancher TerraformAccTest"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "container_format", "bare"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "disk_format", "qcow2"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "schema", "/v2/schemas/image"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "hidden", "true"),
				),
			},
		},
	})
}

func TestAccImagesImage_name(t *testing.T) {
	var image images.Image

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckImage(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckImagesImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageName1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "name", "Rancher TerraformAccTest"),
				),
			},
			{
				Config: testAccImagesImageName2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "name", "TerraformAccTest Rancher"),
				),
			},
		},
	})
}

func TestAccImagesImage_tags(t *testing.T) {
	var image images.Image

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckImage(t)
		},
		ProviderFactories: testAccProviders,
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
		PreCheck: func() {
			testAccPreCheckImage(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckImagesImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image1),
					resource.TestCheckResourceAttrSet(
						"vkcs_images_image.image_1", "properties.os_hash_value"),
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
					resource.TestCheckResourceAttrSet(
						"vkcs_images_image.image_1", "properties.os_hash_value"),
				),
			},
			{
				Config: testAccImagesImageProperties2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image3),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "properties.foo", "bar"),
					resource.TestCheckResourceAttrSet(
						"vkcs_images_image.image_1", "properties.os_hash_value"),
				),
			},
			{
				Config: testAccImagesImageProperties3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image4),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "properties.foo", "baz"),
					resource.TestCheckResourceAttrSet(
						"vkcs_images_image.image_1", "properties.os_hash_value"),
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
					resource.TestCheckResourceAttrSet(
						"vkcs_images_image.image_1", "properties.os_hash_value"),
				),
			},
		},
	})
}

func TestAccImagesImage_webdownload(t *testing.T) {
	var image images.Image

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckImage(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckImagesImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageWebdownload,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageExists("vkcs_images_image.image_1", &image),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "name", "Rancher TerraformAccTest"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "container_format", "bare"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "disk_format", "qcow2"),
					resource.TestCheckResourceAttr(
						"vkcs_images_image.image_1", "schema", "/v2/schemas/image"),
				),
			},
		},
	})
}

func testAccCheckImagesImageDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	imageClient, err := config.ImageV2Client(osRegionName)
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

		config := testAccProvider.Meta().(configer)
		imageClient, err := config.ImageV2Client(osRegionName)
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

		config := testAccProvider.Meta().(configer)
		imageClient, err := config.ImageV2Client(osRegionName)
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

		config := testAccProvider.Meta().(configer)
		imageClient, err := config.ImageV2Client(osRegionName)
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
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"

      timeouts {
        create = "10m"
      }
  }`

const testAccImagesImageBasicWithID = `
  resource "vkcs_images_image" "image_1" {
      name = "Rancher TerraformAccTest"
      image_id = "c1efdf94-9a1a-4401-88b8-d616029d2551"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"

      timeouts {
        create = "10m"
      }
  }`

const testAccImagesImageBasicHidden = `
  resource "vkcs_images_image" "image_1" {
      name = "Rancher TerraformAccTest"
      hidden = true
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"

      timeouts {
        create = "10m"
      }
  }`

const testAccImagesImageName1 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"
  }`

const testAccImagesImageName2 = `
  resource "vkcs_images_image" "image_1" {
      name   = "TerraformAccTest Rancher"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"
  }`

const testAccImagesImageTags1 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"
      tags = ["foo","bar"]
  }`

const testAccImagesImageTags2 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"
      tags = ["foo","bar","baz"]
  }`

const testAccImagesImageTags3 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"
      tags = ["foo","baz"]
  }`

const testAccImagesImageProperties1 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"

      properties = {
        foo = "bar"
        bar = "foo"
      }
  }`

const testAccImagesImageProperties2 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"

      properties = {
        foo = "bar"
      }
  }`

const testAccImagesImageProperties3 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"

      properties = {
        foo = "baz"
      }
  }`

const testAccImagesImageProperties4 = `
  resource "vkcs_images_image" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"

      properties = {
        foo = "baz"
        bar = "foo"
      }
  }`

const testAccImagesImageWebdownload = `
  resource "vkcs_images_image" "image_1" {
      name   = "Rancher TerraformAccTest"
      image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
      container_format = "bare"
      disk_format = "qcow2"
      web_download = true

      timeouts {
        create = "10m"
      }
  }`
