package blockstorage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
)

func TestAccBlockStorageVolume_basic(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckBlockStorageVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageVolumeBasic),
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
				Config: acctest.AccTestRenderConfig(testAccBlockStorageVolumeUpdate),
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
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckBlockStorageVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageVolumeOnlineResize),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_volume.volume_1", "size", "1"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageVolumeOnlineResizeUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_volume.volume_1", "size", "2"),
				),
			},
		},
	})
}

func TestAccBlockStorageVolume_retype(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckBlockStorageVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageVolumeRetype),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_volume.volume_1", "volume_type", "ceph-hdd"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageVolumeRetypeUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_volume.volume_1", "volume_type", "ceph-ssd"),
				),
			},
		},
	})
}

func TestAccBlockStorageVolume_image(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckBlockStorageVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageVolumeImage),
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
	config := acctest.AccTestProvider.Meta().(clients.Config)
	blockStorageClient, err := config.BlockStorageV3Client(acctest.OsRegionName)
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

		config := acctest.AccTestProvider.Meta().(clients.Config)
		blockStorageClient, err := config.BlockStorageV3Client(acctest.OsRegionName)
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

const testAccBlockStorageVolumeBasic = `
resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 1
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}`

const testAccBlockStorageVolumeOnlineResize = `
{{.BaseImage}}

resource "vkcs_compute_instance" "basic" {
  name          = "instance_1"
  flavor_name   = "{{.FlavorName}}"
  image_id      = data.vkcs_images_image.base.id
  network_mode  = "none"
}

resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  description = "test volume"
  size = 1
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}

resource "vkcs_compute_volume_attach" "va_1" {
  instance_id = vkcs_compute_instance.basic.id
  volume_id   = vkcs_blockstorage_volume.volume_1.id
}
`

const testAccBlockStorageVolumeOnlineResizeUpdate = `
{{.BaseImage}}

resource "vkcs_compute_instance" "basic" {
  name          = "instance_1"
  flavor_name   = "{{.FlavorName}}"
  image_id      = data.vkcs_images_image.base.id
  network_mode  = "none"
}

resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  description = "test volume"
  size = 2
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}

resource "vkcs_compute_volume_attach" "va_1" {
  instance_id = vkcs_compute_instance.basic.id
  volume_id   = vkcs_blockstorage_volume.volume_1.id
}
`

const testAccBlockStorageVolumeRetype = `
resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 1
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "ceph-hdd"
}
`

const testAccBlockStorageVolumeRetypeUpdate = `
resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 1
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "ceph-ssd"
}
`

const testAccBlockStorageVolumeUpdate = `
resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1-updated"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 2
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}
`

const testAccBlockStorageVolumeImage = `
{{.BaseImage}}

resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 5
  image_id = data.vkcs_images_image.base.id
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}
`
