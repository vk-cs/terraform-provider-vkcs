package firewall_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccNetworkingSecGroup_importBasic(t *testing.T) {
	resourceName := "vkcs_networking_secgroup.secgroup_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckNetworkingSecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
