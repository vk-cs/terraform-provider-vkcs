package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccLBListener_importBasic(t *testing.T) {
	resourceName := "vkcs_lb_listener.listener_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbListenerConfigBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLBListener_importOctavia(t *testing.T) {
	resourceName := "vkcs_lb_listener.listener_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbListenerConfigOctavia,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
