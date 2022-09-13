package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLBL7Rule_importBasic(t *testing.T) {
	l7ruleResourceName := "vkcs_lb_l7rule.l7rule_1"
	l7policyResourceName := "vkcs_lb_l7policy.l7policy_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBL7RuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccCheckLbL7RuleConfigBasic, map[string]string{"TestAccCheckLbL7RuleConfig": testAccCheckLbL7RuleConfig}),
			},

			{
				ResourceName:      l7ruleResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccLBL7RuleImportID(l7policyResourceName, l7ruleResourceName),
			},
		},
	})
}

func testAccLBL7RuleImportID(l7policyResource, l7ruleResource string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		l7policy, ok := s.RootModule().Resources[l7policyResource]
		if !ok {
			return "", fmt.Errorf("Pool not found: %s", l7policyResource)
		}

		l7rule, ok := s.RootModule().Resources[l7ruleResource]
		if !ok {
			return "", fmt.Errorf("L7Rule not found: %s", l7ruleResource)
		}

		return fmt.Sprintf("%s/%s", l7policy.Primary.ID, l7rule.Primary.ID), nil
	}
}
