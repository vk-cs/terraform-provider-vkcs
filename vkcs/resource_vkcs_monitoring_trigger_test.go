package vkcs

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

const testAccMonitoringTriggerBasic = `
 resource "vkcs_monitoring_channel" "basic" {
   name      = "basic_test"
   channel_type = "email"
   address = "bar@example.com"
 }

 resource "vkcs_monitoring_trigger" "basic" {
 	name = "test"
	namespace = "mcs/test"
	query = "test{} > 10"
	interval = 60
	notification_title = "123"
	notification_channels =  [vkcs_monitoring_channel.basic.id]
 }
`
const testAccMonitoringTriggerUpdate = `
 resource "vkcs_monitoring_channel" "basic" {
   name      = "basic_test"
   channel_type = "email"
   address = "bar@example.com"
 }

 resource "vkcs_monitoring_trigger" "basic" {
 	name = "test"
	namespace = "mcs/test"
	query = "test{} > 10"
	interval = 120
	notification_title = "123"
	notification_channels =  [vkcs_monitoring_channel.basic.id]
 }
`

func testAccCheckMonTriggerExists(n string, ch *TriggerOut) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no id is set")
		}

		config := testAccProvider.Meta().(configer)
		MonClient, err := config.MonitoringV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		found, err := triggerGet(MonClient, config.GetTenantID(), rs.Primary.ID).extract()
		if err != nil {
			return err
		}

		if found.Trigger.Id != rs.Primary.ID {
			return fmt.Errorf("trigger not found")
		}

		*ch = *found

		return nil
	}
}

func TestAccMonitoringTrigger_basic(t *testing.T) {
	var tr TriggerOut

	resource.Test(t, resource.TestCase{
		PreCheck:          func() {},
		ProviderFactories: testAccProviders,
		CheckDestroy:      func(*terraform.State) error { return nil },
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccMonitoringTriggerBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonTriggerExists(
						"vkcs_monitoring_trigger.basic", &tr),
					resource.TestCheckResourceAttrPtr(
						"vkcs_monitoring_trigger.basic", "name", &tr.Trigger.Name),
					resource.TestCheckResourceAttr(
						"vkcs_monitoring_trigger.basic", "interval", "60"),
				),
			},
			{
				Config: testAccRenderConfig(testAccMonitoringTriggerUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonTriggerExists(
						"vkcs_monitoring_trigger.basic", &tr),
					resource.TestCheckResourceAttr(
						"vkcs_monitoring_trigger.basic", "interval", "120"),
				),
			},
		},
	})
}
