package vpnaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccEndpointGroup_importBasic(t *testing.T) {
	resourceName := "vkcs_vpnaas_endpoint_group.group_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
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
