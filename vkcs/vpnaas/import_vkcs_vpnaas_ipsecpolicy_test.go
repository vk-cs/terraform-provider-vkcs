package vpnaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccVPNaaSIPSecPolicy_importBasic(t *testing.T) {
	resourceName := "vkcs_vpnaas_ipsec_policy.policy_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
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
