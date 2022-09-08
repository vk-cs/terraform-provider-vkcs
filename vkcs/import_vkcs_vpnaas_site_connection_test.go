package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSiteConnection_importBasic(t *testing.T) {
	resourceName := "vkcs_vpnaas_site_connection.conn_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSiteConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccSiteConnectionBasic, testAccValues),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
