package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEndpointGroup_importBasic(t *testing.T) {
	resourceName := "vkcs_vpnaas_endpoint_group.group_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckEndpointGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointGroupBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
