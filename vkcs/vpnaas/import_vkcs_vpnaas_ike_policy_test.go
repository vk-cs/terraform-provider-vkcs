package vpnaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccIKEPolicy_importBasic(t *testing.T) {
	resourceName := "vkcs_vpnaas_ike_policy.policy_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
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
