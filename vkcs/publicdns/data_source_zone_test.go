package publicdns_test

import (
	"fmt"
	"testing"

	fm_acctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccPublicDNSZoneDataSource_basic(t *testing.T) {
	zoneName := fmt.Sprintf("vkcs-tf-acctest-%s.com", fm_acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccPublicDNSZoneDataSourceConfigBasic, map[string]string{"ZoneName": zoneName}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicDNSDataSourceID("data.vkcs_publicdns_zone.zone_1"),
					resource.TestCheckResourceAttr("data.vkcs_publicdns_zone.zone_1", "primary_dns", "ns1.mcs.mail.ru"),
					resource.TestCheckResourceAttr("data.vkcs_publicdns_zone.zone_1", "admin_email", "admin@example.com"),
					resource.TestCheckResourceAttr("data.vkcs_publicdns_zone.zone_1", "zone", zoneName),
				),
			},
		},
	})
}

func testAccCheckPublicDNSDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find router data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Router data source ID not set")
		}

		return nil
	}
}

const testAccPublicDNSZoneDataSourceConfigBasic = `
resource "vkcs_publicdns_zone" "zone_1" {
  primary_dns = "ns1.mcs.mail.ru"
  admin_email = "admin@example.com"
  expire = 3600000
  zone = "{{.ZoneName}}"
}

data "vkcs_publicdns_zone" "zone_1" {
  zone = vkcs_publicdns_zone.zone_1.zone
}
`
