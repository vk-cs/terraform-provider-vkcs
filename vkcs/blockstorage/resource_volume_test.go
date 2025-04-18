package blockstorage_test

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	ivolumes "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/blockstorage/v3/volumes"
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
					testAccCheckBlockStorageVolumeMetadata(&volume, map[string]string{"foo": "bar"}),
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
					testAccCheckBlockStorageVolumeMetadata(&volume, map[string]string{"foo": "bar"}),
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

func TestAccBlockStorageVolume_online_retype(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckBlockStorageVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageVolumeOnlineRetype, map[string]string{"VolumeType": "ceph-hdd"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_volume.bootable", "volume_type", "ceph-hdd"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageVolumeOnlineRetype, map[string]string{"VolumeType": "ceph-ssd"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_blockstorage_volume.bootable", "volume_type", "ceph-ssd"),
				),
			},
		},
	})
}

func TestAccBlockStorageVolume_metadata(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckBlockStorageVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageVolumeMetadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageVolumeExists("vkcs_blockstorage_volume.volume_1", &volume),
					testAccCheckBlockStorageVolumeMetadata(&volume, map[string]string{"key1": "val0", "key2": "val2", "key3": "val3"}),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageVolumeMetadataUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageVolumeExists("vkcs_blockstorage_volume.volume_1", &volume),
					resource.TestCheckNoResourceAttr("vkcs_blockstorage_volume.volume_1", "metadata.key2"),
					testAccCheckBlockStorageVolumeMetadata(&volume, map[string]string{"key1": "val1", "key3": "val3", "key4": "val4"}),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccBlockStorageVolumeMetadataDelete),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageVolumeExists("vkcs_blockstorage_volume.volume_1", &volume),
					testAccCheckBlockStorageVolumeEmptyMetadata(&volume),
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

		_, err := ivolumes.Get(blockStorageClient, rs.Primary.ID).Extract()
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

		found, err := ivolumes.Get(blockStorageClient, rs.Primary.ID).Extract()
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

// testAccCheckBlockStorageVolumeMetadata checks that the cloud metadata contains all key-value pairs from the given metadata
func testAccCheckBlockStorageVolumeMetadata(volume *volumes.Volume, metadata map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if volume.Metadata == nil {
			return fmt.Errorf("No metadata")
		}

		for key, value := range metadata {
			cloudVal, exists := volume.Metadata[key]
			if !exists {
				return fmt.Errorf("Metadata not found: %s", key)
			}
			if cloudVal != value {
				return fmt.Errorf("Bad value for key %s: expected: %s, found: %s", key, value, cloudVal)
			}
		}

		return nil
	}
}

func testAccCheckBlockStorageVolumeEmptyMetadata(volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(volume.Metadata) > 0 {
			return fmt.Errorf("Metadata is not empty")
		}

		return nil
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

const testAccBlockStorageVolumeOnlineRetype = `
{{.BaseImage}}

resource "vkcs_blockstorage_volume" "bootable" {
  name              = "bootable-tf-acc"
  size              = 10
  volume_type       = "{{.VolumeType}}"
  image_id          = data.vkcs_images_image.base.id
  availability_zone = "GZ1"
}

resource "vkcs_compute_instance" "basic" {
  name 				= "instance_tf_acc"
  flavor_name       = "{{.FlavorName}}"

  block_device {
    boot_index       = 0
    source_type      = "volume"
    destination_type = "volume"
    uuid             = vkcs_blockstorage_volume.bootable.id
  }

   network_mode  = "none"
}
`

const testAccBlockStorageVolumeMetadata = `
resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  description = "metadata test volume"
  metadata = {
    key1 = "val0"
    key2 = "val2"
    key3 = "val3"
  }
  size = 1
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}`

const testAccBlockStorageVolumeMetadataUpdate = `
resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  description = "metadata test volume"
  metadata = {
    key1 = "val1"
    key3 = "val3"
    key4 = "val4"
  }
  size = 1
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}`

const testAccBlockStorageVolumeMetadataDelete = `
resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  description = "metadata test volume"
  size = 1
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}`
