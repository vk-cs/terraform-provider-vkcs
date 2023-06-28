package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccLBL7Policy_importBasic(t *testing.T) {
	resourceName := "vkcs_lb_l7policy.l7policy_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckLBL7PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccCheckLbL7PolicyConfigBasic, map[string]string{"TestAccCheckLbL7PolicyConfig": testAccCheckLbL7PolicyConfig}),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
