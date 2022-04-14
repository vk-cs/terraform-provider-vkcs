package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLBLoadBalancer_importBasic(t *testing.T) {
	resourceName := "vkcs_lb_loadbalancer.loadbalancer_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckLB(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbLoadBalancerConfigBasic(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
