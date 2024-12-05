package db_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDatabaseDataSourceBackup_basic(t *testing.T) {
	resourceName := "vkcs_db_backup.basic"
	datasourceName := "data.vkcs_db_backup.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDataSourceDatabaseBackupBasic,
					map[string]string{"TestAccDatabaseBackupBasic": acctest.AccTestRenderConfig(testAccDatabaseBackupBasic, map[string]string{"TestAccDatabaseInstanceBasic": acctest.AccTestRenderConfig(testAccDatabaseInstanceBasic)})}),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDatabaseBackupID(datasourceName),
					resource.TestCheckResourceAttrPair(resourceName, "name", datasourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceDatabaseBackupID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find backup data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Backup data source ID not set")
		}

		return nil
	}
}

const testAccDataSourceDatabaseBackupBasic = `
{{.TestAccDatabaseBackupBasic}}

data "vkcs_db_backup" "basic" {
	id = vkcs_db_backup.basic.id
}
`
