package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLBMonitor_importBasic(t *testing.T) {
	resourceName := "vkcs_lb_monitor.monitor_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbMonitorConfigBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
