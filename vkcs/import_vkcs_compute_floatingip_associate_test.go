package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeFloatingIPAssociate_importBasic(t *testing.T) {
	resourceName := "vkcs_compute_floatingip_associate.fip_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckCompute(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeFloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFloatingIPAssociateBasic(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"wait_until_associated",
				},
			},
		},
	})
}
