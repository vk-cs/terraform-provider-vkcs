package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLBL7Policy_importBasic(t *testing.T) {
	resourceName := "vkcs_lb_l7policy.l7policy_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBL7PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccCheckLbL7PolicyConfigBasic, map[string]string{"TestAccCheckLbL7PolicyConfig": testAccCheckLbL7PolicyConfig}),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
