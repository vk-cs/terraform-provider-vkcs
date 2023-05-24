package compute_test

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/compute"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
)

func TestAccComputeVolumeAttach_basic(t *testing.T) {
	var va volumeattach.VolumeAttachment

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckComputeVolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccComputeVolumeAttachBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeVolumeAttachExists("vkcs_compute_volume_attach.va_1", &va),
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

		config := acctest.AccTestProvider.Meta().(clients.Config)
		computeClient, err := config.ComputeV2Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		instanceID, volumeID, err := compute.ComputeVolumeAttachParseID(rs.Primary.ID)
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

const testAccComputeVolumeAttachBasic = `
{{.BaseNetwork}}
{{.BaseImage}}
{{.BaseFlavor}}

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
  security_groups = ["default"]
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
