package firewall_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccNetworkingSecGroupRule_importBasic(t *testing.T) {
	resourceName := "vkcs_networking_secgroup_rule.secgroup_rule_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckNetworkingSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRuleBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
