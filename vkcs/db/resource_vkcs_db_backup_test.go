package db_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/backups"
)

func TestAccDatabaseBackup_basic(t *testing.T) {
	var backup backups.BackupResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseBackupBasic, map[string]string{"TestAccDatabaseInstanceBasic": acctest.AccTestRenderConfig(testAccDatabaseInstanceBasic)}),
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

func testAccCheckDatabaseBackupExists(n string, backup *backups.BackupResp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no id is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		DatabaseClient, err := config.DatabaseV1Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		found, err := backups.Get(DatabaseClient, rs.Primary.ID).Extract()
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
	config := acctest.AccTestProvider.Meta().(clients.Config)

	DatabaseClient, err := config.DatabaseV1Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS database client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_db_backup" {
			continue
		}
		_, err := backups.Get(DatabaseClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("backup still exists")
		}
	}

	return nil
}

const testAccDatabaseBackupBasic = `
{{.TestAccDatabaseInstanceBasic}}

resource "vkcs_db_backup" "basic" {
    name = "basic"
    dbms_id = vkcs_db_instance.basic.id
    description = "basic description"

}
`
