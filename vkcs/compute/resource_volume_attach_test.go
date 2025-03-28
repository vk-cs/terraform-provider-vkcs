package compute_test

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/compute"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	ivolumes "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/blockstorage/v3/volumes"
	ivolumeattach "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/volumeattach"
)

func TestAccComputeVolumeAttach_basic(t *testing.T) {
	var va volumeattach.VolumeAttachment
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeVolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeVolumeAttachBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeVolumeAttachExists("vkcs_compute_volume_attach.va_1", &va),
					testAccCheckBlockStorageVolumeExists("vkcs_blockstorage_volume.volume_1", &volume),

					resource.TestCheckResourceAttr("vkcs_blockstorage_volume.volume_1", "metadata.key1", "val0"),
					resource.TestCheckResourceAttr("vkcs_blockstorage_volume.volume_1", "metadata.key2", "val2"),
					resource.TestCheckResourceAttr("vkcs_blockstorage_volume.volume_1", "metadata.key3", "val3"),

					testAccCheckBlockStorageVolumeMetadata(&volume, map[string]string{"key1": "val0", "key2": "val2", "key3": "val3"}),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeVolumeAttachBasicMetadataUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeVolumeAttachExists("vkcs_compute_volume_attach.va_1", &va),
					testAccCheckBlockStorageVolumeExists("vkcs_blockstorage_volume.volume_1", &volume),

					resource.TestCheckResourceAttr("vkcs_blockstorage_volume.volume_1", "metadata.key1", "val1"),
					resource.TestCheckResourceAttr("vkcs_blockstorage_volume.volume_1", "metadata.key3", "val3"),
					resource.TestCheckResourceAttr("vkcs_blockstorage_volume.volume_1", "metadata.key4", "val4"),
					resource.TestCheckNoResourceAttr("vkcs_blockstorage_volume.volume_1", "metadata.key2"),

					testAccCheckBlockStorageVolumeMetadata(&volume, map[string]string{"key1": "val1", "key3": "val3", "key4": "val4",
						"attached_mode": "rw", "readonly": "False"}),
					testAccCheckBlockStorageVolumeNoMetadata(&volume, []string{"key2"}),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccComputeVolumeAttachBasicMetadataDelete),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeVolumeAttachExists("vkcs_compute_volume_attach.va_1", &va),
					testAccCheckBlockStorageVolumeExists("vkcs_blockstorage_volume.volume_1", &volume),

					testAccCheckBlockStorageVolumeMetadata(&volume, map[string]string{"attached_mode": "rw", "readonly": "False"}),
					testAccCheckBlockStorageVolumeNoMetadata(&volume, []string{"key1", "key2", "key3", "key4"}),
				),
			},
		},
	})
}

func testAccCheckComputeVolumeAttachDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	computeClient, err := config.ComputeV2Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_compute_volume_attach" {
			continue
		}

		instanceID, volumeID, err := compute.ComputeVolumeAttachParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = ivolumeattach.Get(computeClient, instanceID, volumeID).Extract()
		if err == nil {
			return fmt.Errorf("Volume attachment still exists")
		}
	}

	return nil
}

func testAccCheckComputeVolumeAttachExists(n string, va *volumeattach.VolumeAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		computeClient, err := config.ComputeV2Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		instanceID, volumeID, err := compute.ComputeVolumeAttachParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := ivolumeattach.Get(computeClient, instanceID, volumeID).Extract()
		if err != nil {
			return err
		}

		if found.ServerID != instanceID || found.VolumeID != volumeID {
			return fmt.Errorf("VolumeAttach not found")
		}

		*va = *found

		return nil
	}
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

// testAccCheckBlockStorageVolumeMetadata checks that the cloud metadata doesn't contain all keys from the given metadata
func testAccCheckBlockStorageVolumeNoMetadata(volume *volumes.Volume, metadataKeys []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if volume.Metadata == nil {
			return nil
		}

		for _, key := range metadataKeys {
			if cloudVal, exists := volume.Metadata[key]; exists {
				return fmt.Errorf("Unexpected metadata found: %s:%s", key, cloudVal)
			}
		}

		return nil
	}
}

const testAccComputeVolumeAttachBasic = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 1
  metadata = {
    key1 = "val0"
    key2 = "val2"
    key3 = "val3"
  }
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_compute_volume_attach" "va_1" {
  instance_id = vkcs_compute_instance.instance_1.id
  volume_id = vkcs_blockstorage_volume.volume_1.id
}
`

const testAccComputeVolumeAttachBasicMetadataUpdate = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 1
  metadata = {
    key1 = "val1"
    key3 = "val3"
    key4 = "val4"
  }
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_compute_volume_attach" "va_1" {
  instance_id = vkcs_compute_instance.instance_1.id
  volume_id = vkcs_blockstorage_volume.volume_1.id
}
`

const testAccComputeVolumeAttachBasicMetadataDelete = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}
{{.BaseSecurityGroup}}

resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 1
  availability_zone = "{{.AvailabilityZone}}"
  volume_type = "{{.VolumeType}}"
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_router_interface.base"]
  name = "instance_1"
  availability_zone = "{{.AvailabilityZone}}"
  security_group_ids = [data.vkcs_networking_secgroup.default_secgroup.id]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_compute_volume_attach" "va_1" {
  instance_id = vkcs_compute_instance.instance_1.id
  volume_id = vkcs_blockstorage_volume.volume_1.id
}
`
