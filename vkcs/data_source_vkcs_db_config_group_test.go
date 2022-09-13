package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDatabaseDataSourceConfigGroup_basic(t *testing.T) {
	resourceName := "vkcs_db_config_group.basic"
	datasourceName := "data.vkcs_db_config_group.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseConfigGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatabaseConfigGroupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDatabaseConfigGroupID(datasourceName),
					resource.TestCheckResourceAttrPair(resourceName, "name", datasourceName, "name"),
					resource.TestCheckResourceAttr("vkcs_db_config_group.basic", "values.max_connections", "100"),
				),
			},
		},
	})
}

func testAccDataSourceDatabaseConfigGroupID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find config group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Config group data source ID not set")
		}

		return nil
	}
}

var testAccDataSourceDatabaseConfigGroupBasic = fmt.Sprintf(`
%s

data "vkcs_db_config_group" "basic" {
	config_group_id = "${vkcs_db_config_group.basic.id}"
}
`, testAccDatabaseConfigGroupResource)
