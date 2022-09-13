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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseConfigGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccDatabaseConfigGroupBasic),
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
				Config: testAccRenderConfig(testAccDatabaseConfigGroupUpdate),
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

var testAccDatabaseConfigGroupResource = `
resource "vkcs_db_config_group" "basic" {
	name = "basic"
	datastore {
		version = "13"
		type = "postgresql"
	}
	values = {
		max_connections: "100"
	}
}
`

const testAccDatabaseConfigGroupBasic = `
{{.BaseNetwork}}	
{{.BaseFlavor}}

resource "vkcs_db_config_group" "basic" {
	name = "basic"
	datastore {
		version = "13"
		type = "postgresql"
	}
	values = {
		max_connections: "100"
	}
}

resource "vkcs_db_instance" "basic" {
  name             = "basic"
  flavor_id = data.vkcs_compute_flavor.base.id
  size = 8
  volume_type = "{{.VolumeType}}"
  configuration_id = vkcs_db_config_group.basic.id
  datastore {
    version = "13"
    type    = "postgresql"
  }

  network {
    uuid = vkcs_networking_network.base.id
  }
  availability_zone = "{{.AvailabilityZone}}"
  floating_ip_enabled = true

  disk_autoexpand {
    autoexpand = true
    max_disk_size = 1000
  }
  depends_on = [
    vkcs_networking_network.base,
    vkcs_networking_subnet.base
  ]
}
`

const testAccDatabaseConfigGroupUpdate = `
{{.BaseNetwork}}
{{.BaseFlavor}}

resource "vkcs_db_config_group" "basic" {
	name = "basic"
	datastore {
		version = "13"
		type = "postgresql"
	}
	values = {
		max_connections: "200"
	}
}

resource "vkcs_db_instance" "basic" {
  name             = "basic"
  flavor_id = data.vkcs_compute_flavor.base.id
  size = 8
  volume_type = "{{.VolumeType}}"
  configuration_id = vkcs_db_config_group.basic.id
  datastore {
    version = "13"
    type    = "postgresql"
  }

  network {
    uuid = vkcs_networking_network.base.id
  }
  availability_zone = "{{.AvailabilityZone}}"
  floating_ip_enabled = true

  disk_autoexpand {
    autoexpand = true
    max_disk_size = 1000
  }
  depends_on = [
    vkcs_networking_network.base,
    vkcs_networking_subnet.base
  ]
}
`
