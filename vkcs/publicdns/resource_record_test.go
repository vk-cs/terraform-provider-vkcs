package publicdns_test

import (
	"fmt"
	"testing"

	fm_acctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/publicdns/v2/records"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/publicdns/v2/zones"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/publicdns"
)

func TestAccPublicDNSRecord_basic(t *testing.T) {
	var z zones.Zone
	var r map[string]interface{}
	zoneName := fmt.Sprintf("vkcs-tf-acctest-%s.com", fm_acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckPublicDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfigWithZone(testAccPublicDNSRecordConfigBasic, zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicDNSZoneExists("vkcs_publicdns_zone.zone_1", &z),
					testAccCheckPublicDNSRecordExists("vkcs_publicdns_record.record_a", &r),
					resource.TestCheckResourceAttrPtr("vkcs_publicdns_record.record_a", "zone_id", &z.ID),
					resource.TestCheckResourceAttr("vkcs_publicdns_record.record_a", "ip", "8.8.8.8"),
					resource.TestCheckResourceAttr("vkcs_publicdns_record.record_a", "ttl", "60"),
					resource.TestCheckResourceAttr("vkcs_publicdns_record.record_srv", "full_name", "_sip._udp"),
				),
			},
		},
	})
}

func TestAccPublicDNSRecord_update(t *testing.T) {
	zoneName := fmt.Sprintf("vkcs-tf-acctest-%s.com", fm_acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckPublicDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfigWithZone(testAccPublicDNSRecordConfigBasic, zoneName),
			},
			{
				Config: testAccRenderConfigWithZone(testAccPublicDNSRecordConfigUpdate, zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_publicdns_record.record_a", "type", "AAAA"),
					resource.TestCheckResourceAttr("vkcs_publicdns_record.record_a", "ip", "2001:db8:aa10:1::fb"),
					resource.TestCheckResourceAttr("vkcs_publicdns_record.record_a", "ttl", "86400"),
					resource.TestCheckResourceAttr("vkcs_publicdns_record.record_srv", "full_name", "_sip._tcp.test"),
				),
			},
		},
	})
}

func testAccCheckPublicDNSRecordDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	client, err := config.PublicDNSV2Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS public DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_publicdns_record" {
			continue
		}

		zoneID := rs.Primary.Attributes["zone_id"]
		recordType := rs.Primary.Attributes["type"]
		res := records.Get(client, zoneID, rs.Primary.ID, recordType)
		_, err := publicdns.PublicDNSRecordExtract(res, recordType)
		if err == nil {
			return fmt.Errorf("Record still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckPublicDNSRecordExists(n string, r *map[string]interface{}) resource.TestCheckFunc {
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

		zoneID := rs.Primary.Attributes["zone_id"]
		recordType := rs.Primary.Attributes["type"]
		res := records.Get(client, zoneID, rs.Primary.ID, recordType)

		record, err := publicdns.PublicDNSRecordExtract(res, recordType)
		if err != nil {
			return err
		}

		found, err := util.StructToMap(record)
		if err != nil {
			return err
		}

		if found["uuid"] != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*r = found

		return nil
	}
}

const testAccPublicDNSRecordConfigBasic = `
resource "vkcs_publicdns_zone" "zone_1" {
  primary_dns = "ns1.mcs.mail.ru"
  admin_email = "admin@example.com"
  expire = 3600000
  zone = "{{.ZoneName}}"
}

resource "vkcs_publicdns_record" "record_a" {
  zone_id = vkcs_publicdns_zone.zone_1.id
  type = "A"
  name = "test"
  ip = "8.8.8.8"
  ttl = 60
}

resource "vkcs_publicdns_record" "record_aaaa" {
  zone_id = vkcs_publicdns_zone.zone_1.id
  type = "AAAA"
  name = "test"
  ip = "2001:0DB8:AA10:0001:0000:0000:0000:00FB"
  ttl = 86400
}

resource "vkcs_publicdns_record" "record_cname" {
  zone_id = vkcs_publicdns_zone.zone_1.id
  type = "CNAME"
  name = "alias"
  content = "example.com"
}

resource "vkcs_publicdns_record" "record_mx" {
  zone_id = vkcs_publicdns_zone.zone_1.id
  type = "MX"
  name = "@"
  priority = 10
  content = "mx.example.com"
}

resource "vkcs_publicdns_record" "record_ns" {
  zone_id = vkcs_publicdns_zone.zone_1.id
  type = "NS"
  name = "@"
  content = "ns1.example.com"
}

resource "vkcs_publicdns_record" "record_srv" {
  zone_id = vkcs_publicdns_zone.zone_1.id
  type = "SRV"
  service = "_sip"
  proto = "_udp"
  priority = 10
  name = ""
  weight = 5
  host = "example.com"
  port = 5060
  ttl = 86400
}

resource "vkcs_publicdns_record" "record_txt" {
  zone_id = vkcs_publicdns_zone.zone_1.id
  type = "TXT"
  name = ""
  content = "Example"
}
`

const testAccPublicDNSRecordConfigUpdate = `
resource "vkcs_publicdns_zone" "zone_1" {
  primary_dns = "ns1.mcs.mail.ru"
  admin_email = "admin@example.com"
  expire = 3600000
  zone = "{{.ZoneName}}"
}

resource "vkcs_publicdns_record" "record_a" {
  zone_id = vkcs_publicdns_zone.zone_1.id
  type = "AAAA"
  name = "test"
  ip = "2001:0DB8:AA10:0001:0000:0000:0000:00FB"
  ttl = 86400
}

resource "vkcs_publicdns_record" "record_aaaa" {
  zone_id = vkcs_publicdns_zone.zone_1.id
  type = "A"
  name = "test"
  ip = "8.8.8.8"
  ttl = 86400
}

resource "vkcs_publicdns_record" "record_cname" {
  zone_id = vkcs_publicdns_zone.zone_1.id
  type = "CNAME"
  name = ""
  content = "new-example.com"
}

resource "vkcs_publicdns_record" "record_mx" {
  zone_id = vkcs_publicdns_zone.zone_1.id
  type = "MX"
  name = "@"
  priority = 20
  content = "mx.example.com"
}

resource "vkcs_publicdns_record" "record_ns" {
  zone_id = vkcs_publicdns_zone.zone_1.id
  type = "NS"
  name = "@"
  content = "ns2.example.com"
}

resource "vkcs_publicdns_record" "record_srv" {
  zone_id = vkcs_publicdns_zone.zone_1.id
  type = "SRV"
  service = "_sip"
  proto = "_tcp"
  priority = 20
  weight = 1
  name = "test"
  host = "new-example.com"
  port = 5061
  ttl = 3600
}

resource "vkcs_publicdns_record" "record_txt" {
  zone_id = vkcs_publicdns_zone.zone_1.id
  type = "TXT"
  name = "test"
  content = "New example"
}
`
