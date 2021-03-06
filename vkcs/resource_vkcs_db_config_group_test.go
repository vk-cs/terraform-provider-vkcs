package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDatabaseConfigGroup_basic(t *testing.T) {
	var configGroup dbConfigGroupResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDatabase(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseConfigGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseConfigGroupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseConfigGroupExists(
						"vkcs_db_config_group.basic", &configGroup),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_config_group.basic", "name", &configGroup.Name),
					resource.TestCheckResourceAttr(
						"vkcs_db_config_group.basic", "values.max_connections", "100"),
				),
			},
			{
				Config: testAccDatabaseConfigGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseConfigGroupExists(
						"vkcs_db_config_group.basic", &configGroup),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_config_group.basic", "name", &configGroup.Name),
					resource.TestCheckResourceAttr(
						"vkcs_db_config_group.basic", "values.max_connections", "200"),
				),
			},
		},
	})
}

func testAccCheckDatabaseConfigGroupExists(n string, configGroup *dbConfigGroupResp) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating VKCS database client: %s", err)
		}

		found, err := dbConfigGroupGet(DatabaseClient, rs.Primary.ID).extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("instance not found")
		}

		*configGroup = *found

		return nil
	}
}

func testAccCheckDatabaseConfigGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)

	DatabaseClient, err := config.DatabaseV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS database client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_db_config_group" {
			continue
		}
		_, err := instanceGet(DatabaseClient, rs.Primary.ID).extract()
		if err == nil {
			return fmt.Errorf("config group still exists")
		}
	}

	return nil
}

var testAccDatabaseConfigGroupResource = fmt.Sprintf(
	`
resource "vkcs_db_config_group" "basic" {
	name = "basic"
	datastore {
		version = "%s"
		type = "%s"
	}
	values = {
		max_connections: "100"
	}
}
`, osDBDatastoreVersion, osDBDatastoreType)

var testAccDatabaseConfigGroupBasic = fmt.Sprintf(`
%s

resource "vkcs_db_instance" "basic" {
  name             = "basic"
  flavor_id = "%s"
  size = 8
  volume_type = "ms1"
  configuration_id = vkcs_db_config_group.basic.id
  datastore {
    version = "%s"
    type    = "%s"
  }

  network {
    uuid = "%s"
  }
  availability_zone = "MS1"
  floating_ip_enabled = true
  keypair = "%s"

  disk_autoexpand {
    autoexpand = true
    max_disk_size = 1000
  }

}
`, testAccDatabaseConfigGroupResource, osFlavorID, osDBDatastoreVersion, osDBDatastoreType, osNetworkID, osKeypairName)

var testAccDatabaseConfigGroupUpdate = fmt.Sprintf(`
resource "vkcs_db_config_group" "basic" {
	name = "basic"
	datastore {
		version = "%s"
		type = "%s"
	}
	values = {
		max_connections: "200"
	}
}

resource "vkcs_db_instance" "basic" {
  name             = "basic"
  flavor_id = "%s"
  size = 8
  volume_type = "ms1"
  configuration_id = vkcs_db_config_group.basic.id
  datastore {
    version = "%s"
    type    = "%s"
  }

  network {
    uuid = "%s"
  }
  availability_zone = "MS1"
  floating_ip_enabled = true
  keypair = "%s"

  disk_autoexpand {
    autoexpand = true
    max_disk_size = 1000
  }

}
`, osDBDatastoreVersion, osDBDatastoreType, osFlavorID, osDBDatastoreVersion, osDBDatastoreType, osNetworkID, osKeypairName)
