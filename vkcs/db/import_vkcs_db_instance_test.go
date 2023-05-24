package db_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDatabaseInstance_importBasic(t *testing.T) {
	resourceName := "vkcs_db_instance.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseInstanceBasic),
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
