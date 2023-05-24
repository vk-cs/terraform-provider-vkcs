package db_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	cg "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/config_groups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/instances"
)

func TestAccDatabaseConfigGroup_basic(t *testing.T) {
	var configGroup cg.ConfigGroupResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseConfigGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseConfigGroupBasic),
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
				Config: acctest.AccTestRenderConfig(testAccDatabaseConfigGroupUpdate),
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

func testAccCheckDatabaseConfigGroupExists(n string, configGroup *cg.ConfigGroupResp) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating VKCS database client: %s", err)
		}

		found, err := cg.Get(DatabaseClient, rs.Primary.ID).Extract()
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
	config := acctest.AccTestProvider.Meta().(clients.Config)

	DatabaseClient, err := config.DatabaseV1Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS database client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_db_config_group" {
			continue
		}
		_, err := instances.Get(DatabaseClient, rs.Primary.ID).Extract()
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
  depends_on = [vkcs_networking_router_interface.base]
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
  depends_on = [vkcs_networking_router_interface.base]
}
`
