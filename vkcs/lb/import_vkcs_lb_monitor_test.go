package lb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccLBMonitor_importBasic(t *testing.T) {
	resourceName := "vkcs_lb_monitor.monitor_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
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
