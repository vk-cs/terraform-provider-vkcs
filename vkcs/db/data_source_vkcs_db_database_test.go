package db_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDatabaseDataSourceDatabase_basic(t *testing.T) {
	resourceName := "vkcs_db_database.basic"
	datasourceName := "data.vkcs_db_database.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDataSourceDatabaseDatabaseBasic, map[string]string{"TestAccDatabaseDatabaseBasic": acctest.AccTestRenderConfig(testAccDatabaseDatabaseBasic)}),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDatabaseDatabaseID(datasourceName),
					resource.TestCheckResourceAttrPair(resourceName, "name", datasourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceDatabaseDatabaseID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find database data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Database data source ID not set")
		}

		return nil
	}
}

const testAccDataSourceDatabaseDatabaseBasic = `
{{.TestAccDatabaseDatabaseBasic}}

data "vkcs_db_database" "basic" {
	id = vkcs_db_database.basic.id
}
`
