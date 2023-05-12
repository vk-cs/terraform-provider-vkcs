package db_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDatabaseConfigGroup_importBasic(t *testing.T) {
	resourceName := "vkcs_db_config_group.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseConfigGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseConfigGroupResource,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
