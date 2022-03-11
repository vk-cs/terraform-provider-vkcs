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
		PreCheck:          func() { testAccPreCheckCompute(t) },
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
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
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
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
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
resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 1
  availability_zone = "nova"
  volume_type = "%s"
}

resource "vkcs_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "vkcs_compute_volume_attach" "va_1" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  volume_id = "${vkcs_blockstorage_volume.volume_1.id}"
}
`, osVolumeType, osNetworkID)
}
