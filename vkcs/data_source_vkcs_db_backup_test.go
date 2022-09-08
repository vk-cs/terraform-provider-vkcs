package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDatabaseDataSourceBackup_basic(t *testing.T) {
	resourceName := "vkcs_db_backup.basic"
	datasourceName := "data.vkcs_db_backup.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatabaseBackupBasic,
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

var testAccDataSourceDatabaseBackupBasic = fmt.Sprintf(`
%s

data "vkcs_db_backup" "basic" {
	backup_id = "${vkcs_db_backup.basic.id}"
}
`, testAccDatabaseBackupBasic)
