package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIPSecPolicy_importBasic(t *testing.T) {
	resourceName := "vkcs_vpnaas_ipsec_policy.policy_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPSecPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
