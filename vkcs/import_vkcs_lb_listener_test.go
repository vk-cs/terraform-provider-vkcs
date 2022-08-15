package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLBListener_importBasic(t *testing.T) {
	resourceName := "vkcs_lb_listener.listener_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
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
		ProviderFactories: testAccProviders,
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
