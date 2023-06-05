package firewall_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccFirewallSecGroupRule_importBasic(t *testing.T) {
	resourceName := "vkcs_networking_secgroup_rule.secgroup_rule_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccFirewallCheckSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallSecGroupRuleBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
