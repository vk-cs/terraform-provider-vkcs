package vkcs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/db/v1/databases"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDatabaseDatabase_basic(t *testing.T) {
	var database databases.Database
	var instance instanceResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDatabase(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseDatabaseBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseInstanceExists(
						"vkcs_db_instance.basic", &instance),
					testAccCheckDatabaseDatabaseExists(
						"vkcs_db_database.basic", &instance, &database),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_database.basic", "name", &database.Name),
				),
			},
		},
	})
}

func testAccCheckDatabaseDatabaseExists(n string, instance *instanceResp, database *databases.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		parts := strings.SplitN(rs.Primary.ID, "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Malformed database name: %s", rs.Primary.ID)
		}

		config := testAccProvider.Meta().(configer)
		DatabaseClient, err := config.DatabaseV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating cloud database client: %s", err)
		}

		pages, err := databaseList(DatabaseClient, instance.ID, "instance").AllPages()
		if err != nil {
			return fmt.Errorf("Unable to retrieve databases: %s", err)
		}

		allDatabases, err := databases.ExtractDBs(pages)
		if err != nil {
			return fmt.Errorf("Unable to extract databases: %s", err)
		}

		for _, db := range allDatabases {
			if db.Name == parts[1] {
				*database = db
				return nil
			}
		}

		return fmt.Errorf("Database %s does not exist", n)
	}
}

func testAccCheckDatabaseDatabaseDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)

	DatabaseClient, err := config.DatabaseV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating cloud database client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_db_database" {
			continue
		}

		parts := strings.SplitN(rs.Primary.ID, "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Malformed database name: %s", rs.Primary.ID)
		}

		pages, err := databaseList(DatabaseClient, parts[0], "instance").AllPages()
		if err != nil {
			return nil
		}

		allDatabase, err := databases.ExtractDBs(pages)
		if err != nil {
			return fmt.Errorf("Unable to extract databases: %s", err)
		}

		var exists bool
		for _, db := range allDatabase {
			if db.Name == parts[1] {
				exists = true
			}
		}

		if exists {
			return fmt.Errorf("Database still exists")
		}
	}

	return nil
}

var testAccDatabaseDatabaseBasic = fmt.Sprintf(`
resource "vkcs_db_instance" "basic" {
  name = "basic"
  flavor_id = "%s"
  size = 10
  volume_type = "ms1"

  datastore {
    version = "%s"
    type    = "%s"
  }

  network {
    uuid = "%s"
  }
}

resource "vkcs_db_database" "basic" {
  name        = "basic"
  dbms_id = "${vkcs_db_instance.basic.id}"
}
`, osFlavorID, osDBDatastoreVersion, osDBDatastoreType, osNetworkID)
