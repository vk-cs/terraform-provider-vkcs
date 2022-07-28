package vkcs

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeVolumeAttach_basic(t *testing.T) {
	var va volumeattach.VolumeAttachment

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeVolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVolumeAttachBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeVolumeAttachExists("vkcs_compute_volume_attach.va_1", &va),
				),
			},
		},
	})
}

func testAccCheckComputeVolumeAttachDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	computeClient, err := config.ComputeV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_compute_volume_attach" {
			continue
		}

		instanceID, volumeID, err := computeVolumeAttachParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = volumeattach.Get(computeClient, instanceID, volumeID).Extract()
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

		config := testAccProvider.Meta().(configer)
		computeClient, err := config.ComputeV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		instanceID, volumeID, err := computeVolumeAttachParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := volumeattach.Get(computeClient, instanceID, volumeID).Extract()
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

func testAccComputeVolumeAttachBasic() string {
	return fmt.Sprintf(`
%s

%s

%s

resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 1
  availability_zone = "GZ1"
  volume_type = "ceph-ssd"
}

resource "vkcs_compute_instance" "instance_1" {
  depends_on = ["vkcs_networking_subnet.base"]
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = vkcs_networking_network.base.id
  }
  image_id = data.vkcs_images_image.base.id
  flavor_id = data.vkcs_compute_flavor.base.id
}

resource "vkcs_compute_volume_attach" "va_1" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  volume_id = "${vkcs_blockstorage_volume.volume_1.id}"
}
`, testAccBaseFlavor, testAccBaseImage, testAccBaseNetwork)
}
