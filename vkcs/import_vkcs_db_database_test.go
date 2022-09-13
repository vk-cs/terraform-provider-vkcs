package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatabaseDatabase_importBasic(t *testing.T) {
	resourceName := "vkcs_db_database.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccDatabaseDatabaseBasic, testAccValues),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
