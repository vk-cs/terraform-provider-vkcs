package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeVolumeAttach_importBasic(t *testing.T) {
	resourceName := "vkcs_compute_volume_attach.va_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeVolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVolumeAttachBasic(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
