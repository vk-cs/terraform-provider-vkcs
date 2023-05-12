package publicdns_test

import (
	"fmt"
	"testing"

	sdk_acctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccPublicDNSZone_importBasic(t *testing.T) {
	resourceName := "vkcs_publicdns_zone.zone_1"
	zoneName := fmt.Sprintf("vkcs-tf-acctest-%s.com", sdk_acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckPublicDNSZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfigWithZone(testAccPublicDNSZoneConfigBasic, zoneName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
