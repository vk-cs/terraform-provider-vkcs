package publicdns_test

import (
	"fmt"
	"testing"

	fm_acctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccPublicDNSZone_importBasic(t *testing.T) {
	resourceName := "vkcs_publicdns_zone.zone_1"
	zoneName := fmt.Sprintf("vkcs-tf-acctest-%s.com", fm_acctest.RandString(5))

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
