package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccLBLoadBalancer_importBasic(t *testing.T) {
	resourceName := "vkcs_lb_loadbalancer.loadbalancer_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckLBLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbLoadBalancerConfigBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
