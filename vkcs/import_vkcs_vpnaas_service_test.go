package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccService_importBasic(t *testing.T) {
	resourceName := "vkcs_vpnaas_service.service_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccServiceBasic),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
