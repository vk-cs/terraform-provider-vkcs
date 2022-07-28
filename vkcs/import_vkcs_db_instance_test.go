package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatabaseInstance_importBasic(t *testing.T) {
	resourceName := "vkcs_db_instance.basic"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseInstanceBasic,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"volume_type", "network", "keypair", "availability_zone", "floating_ip_enabled", "disk_autoexpand"},
			},
		},
	})
}
