package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIKEPolicy_importBasic(t *testing.T) {
	resourceName := "vkcs_vpnaas_ike_policy.policy_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIKEPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
