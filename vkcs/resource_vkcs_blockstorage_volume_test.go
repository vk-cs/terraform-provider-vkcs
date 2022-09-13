package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
)

func TestAccBlockStorageVolume_basic(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageVolumeBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageVolumeExists("vkcs_blockstorage_volume.volume_1", &volume),
					testAccCheckBlockStorageVolumeMetadata(&volume, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_volume.volume_1", "name", "volume_1"),
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_volume.volume_1", "size", "1"),
				),
			},
			{
				Config: testAccBlockStorageVolumeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageVolumeExists("vkcs_blockstorage_volume.volume_1", &volume),
					testAccCheckBlockStorageVolumeMetadata(&volume, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_volume.volume_1", "name", "volume_1-updated"),
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_volume.volume_1", "size", "2"),
				),
			},
		},
	})
}

func TestAccBlockStorageVolume_online_resize(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageVolumeOnlineResize(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_volume.volume_1", "size", "1"),
				),
			},
			{
				Config: testAccBlockStorageVolumeOnlineResizeUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_volume.volume_1", "size", "2"),
				),
			},
		},
	})
}

func TestAccBlockStorageVolume_image(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageVolumeImage(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageVolumeExists("vkcs_blockstorage_volume.volume_1", &volume),
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_volume.volume_1", "name", "volume_1"),
				),
			},
		},
	})
}

func testAccCheckBlockStorageVolumeDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	blockStorageClient, err := config.BlockStorageV3Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS block storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_blockstorage_volume" {
			continue
		}

		_, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Volume still exists")
		}
	}

	return nil
}

func testAccCheckBlockStorageVolumeExists(n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(configer)
		blockStorageClient, err := config.BlockStorageV3Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS block storage client: %s", err)
		}

		found, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Volume not found")
		}

		*volume = *found

		return nil
	}
}

func testAccCheckBlockStorageVolumeMetadata(
	volume *volumes.Volume, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if volume.Metadata == nil {
			return fmt.Errorf("No metadata")
		}

		for key, value := range volume.Metadata {
			if k != key {
				continue
			}

			if v == value {
				return nil
			}

			return fmt.Errorf("Bad value for %s: %s", k, value)
		}

		return fmt.Errorf("Metadata not found: %s", k)
	}
}

const testAccBlockStorageVolumeBasic string = `
resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 1
  availability_zone = "GZ1"
  volume_type = "ceph-ssd"
}`

func testAccBlockStorageVolumeOnlineResize() string {
	return fmt.Sprintf(`
%s

%s

resource "vkcs_compute_instance" "basic" {
  name          = "instance_1"
  flavor_name   = data.vkcs_compute_flavor.base.name
  image_id      = data.vkcs_images_image.base.id
  network_mode  = "none"
}

resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  description = "test volume"
  size = 1
  availability_zone = "GZ1"
  volume_type = "ceph-ssd"
}

resource "vkcs_compute_volume_attach" "va_1" {
  instance_id = "${vkcs_compute_instance.basic.id}"
  volume_id   = "${vkcs_blockstorage_volume.volume_1.id}"
}
`, testAccBaseFlavor, testAccBaseImage)
}

func testAccBlockStorageVolumeOnlineResizeUpdate() string {
	return fmt.Sprintf(`
%s

%s

resource "vkcs_compute_instance" "basic" {
  name            = "instance_1"
  flavor_name     = data.vkcs_compute_flavor.base.name
  image_id      = data.vkcs_images_image.base.id
  network_mode = "none"
}

resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  description = "test volume"
  size = 2
  availability_zone = "GZ1"
  volume_type = "ceph-ssd"
}

resource "vkcs_compute_volume_attach" "va_1" {
  instance_id = "${vkcs_compute_instance.basic.id}"
  volume_id   = "${vkcs_blockstorage_volume.volume_1.id}"
}
`, testAccBaseFlavor, testAccBaseImage)
}

const testAccBlockStorageVolumeUpdate = `
resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1-updated"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 2
  availability_zone = "GZ1"
  volume_type = "ceph-ssd"
}
`

func testAccBlockStorageVolumeImage() string {
	return fmt.Sprintf(`
%s

resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 5
  image_id = data.vkcs_images_image.base.id
  availability_zone = "GZ1"
  volume_type = "ceph-ssd"
}
`, testAccBaseImage)
}
