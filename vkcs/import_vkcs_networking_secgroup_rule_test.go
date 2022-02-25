package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingSecGroupRule_importBasic(t *testing.T) {
	resourceName := "vkcs_networking_secgroup_rule.secgroup_rule_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckNetworking(t) },
		ProviderFactories: testAccProviders,
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
