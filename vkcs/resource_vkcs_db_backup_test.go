package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDatabaseBackup_basic(t *testing.T) {
	var backup dbBackupResp

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseBackupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseBackupExists(
						"vkcs_db_backup.basic", &backup),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_backup.basic", "name", &backup.Name),
				),
			},
		},
	})
}

func testAccCheckDatabaseBackupExists(n string, backup *dbBackupResp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no id is set")
		}

		config := testAccProvider.Meta().(configer)
		DatabaseClient, err := config.DatabaseV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		found, err := dbBackupGet(DatabaseClient, rs.Primary.ID).extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("backup not found")
		}

		*backup = *found

		return nil
	}
}

func testAccCheckDatabaseBackupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)

	DatabaseClient, err := config.DatabaseV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS database client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_db_backup" {
			continue
		}
		_, err := dbBackupGet(DatabaseClient, rs.Primary.ID).extract()
		if err == nil {
			return fmt.Errorf("backup still exists")
		}
	}

	return nil
}

var testAccDatabaseBackupBasic = fmt.Sprintf(`
%s

resource "vkcs_db_backup" "basic" {
    name = "basic"
    dbms_id = vkcs_db_instance.basic.id
    description = "basic description"

}
`, testAccDatabaseInstanceBasic)
