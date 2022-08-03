package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSFSSecurityService_importBasic(t *testing.T) {
	resourceName := "vkcs_sharedfilesystem_securityservice.securityservice_1"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSSecurityServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSSecurityServiceConfigBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}
