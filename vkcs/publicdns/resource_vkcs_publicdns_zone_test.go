package publicdns_test

import (
	"fmt"
	"testing"

	fm_acctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/publicdns/v2/zones"
)

func TestAccPublicDNSZone_basic(t *testing.T) {
	var z zones.Zone
	zoneName := fmt.Sprintf("vkcs-tf-acctest-%s.com", fm_acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckPublicDNSZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfigWithZone(testAccPublicDNSZoneConfigBasic, zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicDNSZoneExists("vkcs_publicdns_zone.zone_1", &z),
					resource.TestCheckResourceAttr("vkcs_publicdns_zone.zone_1", "primary_dns", "ns1.mcs.mail.ru"),
					resource.TestCheckResourceAttr("vkcs_publicdns_zone.zone_1", "zone", zoneName),
					resource.TestCheckResourceAttrPtr("vkcs_publicdns_zone.zone_1", "status", &z.Status),
				),
			},
			{
				Config: testAccRenderConfigWithZone(testAccPublicDNSZoneConfigUpdate, zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_publicdns_zone.zone_1", "primary_dns", "ns2.mcs.mail.ru"),
					resource.TestCheckResourceAttr("vkcs_publicdns_zone.zone_1", "expire", "7200000"),
					resource.TestCheckResourceAttr("vkcs_publicdns_zone.zone_1", "admin_email", "new-admin@example.com"),
				),
			},
		},
	})
}

func testAccCheckPublicDNSZoneDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	client, err := config.PublicDNSV2Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS public DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_publicdns_zone" {
			continue
		}

		_, err := zones.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Zone still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPublicDNSZoneExists(n string, z *zones.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		client, err := config.PublicDNSV2Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS public DNS client: %s", err)
		}

		found, err := zones.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Zone not found")
		}

		*z = *found

		return nil
	}
}

func testAccRenderConfigWithZone(testConfig string, zoneName string) string {
	return acctest.AccTestRenderConfig(testConfig, map[string]string{"ZoneName": zoneName})
}

const testAccPublicDNSZoneConfigBasic = `
resource "vkcs_publicdns_zone" "zone_1" {
  primary_dns = "ns1.mcs.mail.ru"
  admin_email = "admin@example.com"
  expire = 3600000
  zone = "{{.ZoneName}}"
}
`

const testAccPublicDNSZoneConfigUpdate = `
resource "vkcs_publicdns_zone" "zone_1" {
  primary_dns = "ns2.mcs.mail.ru"
  admin_email = "new-admin@example.com"
  expire = 7200000
  zone = "{{.ZoneName}}"
}
`
