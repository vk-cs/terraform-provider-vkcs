package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLBPool_importBasic(t *testing.T) {
	resourceName := "vkcs_lb_pool.pool_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbPoolConfigBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
