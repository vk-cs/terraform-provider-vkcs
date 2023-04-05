package vkcs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPublicDNSRecord_importBasic(t *testing.T) {
	resourceName := "vkcs_publicdns_record.record_1"
	zoneName := fmt.Sprintf("vkcs-tf-acctest-%s.com", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckPublicDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfigWithZone(testAccPublicDNSRecordConfigBasic, zoneName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					for _, rs := range s.RootModule().Resources {
						if rs.Type != "vkcs_publicdns_record" {
							continue
						}

						zoneID := rs.Primary.Attributes["zone_id"]
						recordType := strings.ToLower(rs.Primary.Attributes["type"])
						recordID := rs.Primary.ID

						return fmt.Sprintf("%s/%s/%s", zoneID, recordType, recordID), nil
					}
					return "", fmt.Errorf("Record not found")
				},
			},
		},
	})
}
