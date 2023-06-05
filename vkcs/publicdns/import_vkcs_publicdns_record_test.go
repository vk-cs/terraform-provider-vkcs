package publicdns_test

import (
	"fmt"
	"strings"
	"testing"

	fm_acctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccPublicDNSRecord_importBasic(t *testing.T) {
	resourceName := "vkcs_publicdns_record.record_a"
	zoneName := fmt.Sprintf("vkcs-tf-acctest-%s.com", fm_acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
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
