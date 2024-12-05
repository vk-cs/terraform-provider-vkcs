package db_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDatabaseDataSourceConfigGroup_basic(t *testing.T) {
	resourceName := "vkcs_db_config_group.basic"
	datasourceName := "data.vkcs_db_config_group.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDataSourceDatabaseConfigGroupBasic, map[string]string{"TestAccDatabaseConfigGroupResource": testAccDatabaseConfigGroupResource}),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDatabaseConfigGroupID(datasourceName),
					resource.TestCheckResourceAttrPair(resourceName, "name", datasourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "values.max_connections", "100"),
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

const testAccDataSourceDatabaseConfigGroupBasic = `
{{.TestAccDatabaseConfigGroupResource}}

data "vkcs_db_config_group" "basic" {
	id = vkcs_db_config_group.basic.id
}
`
