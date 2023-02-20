package vkcs

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccMonitoringChannelBasic = `
 resource "vkcs_monitoring_channel" "basic" {
   name      = "basic_test"
   channel_type = "email"
   address = "foo@example.com"
 }
`
const testAccMonitoringChannelUpdate = `
 resource "vkcs_monitoring_channel" "basic" {
   name      = "basic_test"
   channel_type = "email"
   address = "bar@example.com"
 }
`

func testAccCheckMonChannelExists(n string, ch *ChannelOut) resource.TestCheckFunc {
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

		found, err := channelGet(MonClient, config.GetTenantID(), rs.Primary.ID).extract()
		if err != nil {
			return err
		}

		if found.Channel.Id != rs.Primary.ID {
			return fmt.Errorf("channel not found")
		}

		*ch = *found

		return nil
	}
}

func TestAccMonitoringChannel_basic(t *testing.T) {
	var ch ChannelOut

	resource.Test(t, resource.TestCase{
		PreCheck:          func() {},
		ProviderFactories: testAccProviders,
		CheckDestroy:      func(*terraform.State) error { return nil },
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccMonitoringChannelBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonChannelExists(
						"vkcs_monitoring_channel.basic", &ch),
					resource.TestCheckResourceAttrPtr(
						"vkcs_monitoring_channel.basic", "name", &ch.Channel.Name),
				),
			},
			{
				Config: testAccRenderConfig(testAccMonitoringChannelUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonChannelExists(
						"vkcs_monitoring_channel.basic", &ch),
					resource.TestCheckResourceAttr(
						"vkcs_monitoring_channel.basic", "address", "bar@example.com"),
				),
			},
		},
	})
}
