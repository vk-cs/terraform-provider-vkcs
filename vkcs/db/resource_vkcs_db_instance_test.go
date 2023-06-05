package db_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/instances"
)

func TestAccDatabaseInstance_basic_big(t *testing.T) {
	var instance instances.InstanceResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseInstanceBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseInstanceExists(
						"vkcs_db_instance.basic", &instance),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_instance.basic", "name", &instance.Name),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseInstanceUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseInstanceExists(
						"vkcs_db_instance.basic", &instance),
					resource.TestCheckResourceAttr(
						"vkcs_db_instance.basic", "size", "9"),
				),
			},
		},
	})
}

func TestAccDatabaseInstance_rootUser(t *testing.T) {
	var instance instances.InstanceResp
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseInstanceRootUser),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseInstanceExists(
						"vkcs_db_instance.basic", &instance),
					testAccCheckDatabaseRootUserExists(
						"vkcs_db_instance.basic", &instance),
				),
			},
		},
	})
}

func TestAccDatabaseInstance_wal(t *testing.T) {
	var instance instances.InstanceResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseInstanceWal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseInstanceExists(
						"vkcs_db_instance.basic", &instance),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_instance.basic", "name", &instance.Name),
				),
			},
		},
	})
}

func TestAccDatabaseInstance_wal_no_update(t *testing.T) {
	var instance instances.InstanceResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseInstanceWal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseInstanceExists(
						"vkcs_db_instance.basic", &instance),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_instance.basic", "name", &instance.Name),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseInstanceWal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseInstanceExists(
						"vkcs_db_instance.basic", &instance),
				),
			},
		},
	})
}

func testAccCheckDatabaseInstanceExists(n string, instance *instances.InstanceResp) resource.TestCheckFunc {
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

		found, err := instances.Get(DatabaseClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("instance not found")
		}

		*instance = *found

		return nil
	}
}

func testAccCheckDatabaseInstanceDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)

	DatabaseClient, err := config.DatabaseV1Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS database client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_db_instance" {
			continue
		}
		_, err := instances.Get(DatabaseClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("instance still exists")
		}
	}

	return nil
}

func testAccCheckDatabaseRootUserExists(n string, instance *instances.InstanceResp) resource.TestCheckFunc {

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
			return fmt.Errorf("error creating cloud database client: %s", err)
		}

		isRootEnabledResult := instances.RootUserGet(DatabaseClient, rs.Primary.ID)
		isRootEnabled, err := isRootEnabledResult.Extract()
		if err != nil {
			return fmt.Errorf("error checking if root user is enabled for instance: %s: %s", rs.Primary.ID, err)
		}

		if isRootEnabled {
			return nil
		}

		return fmt.Errorf("root user %s does not exist", n)
	}
}

const testAccDatabaseInstanceBasic = `
{{.BaseNetwork}}
{{.BaseFlavor}}

resource "vkcs_db_instance" "basic" {
  name             = "basic"
  flavor_id = data.vkcs_compute_flavor.base.id
  size = 8
  volume_type = "{{.VolumeType}}"

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

const testAccDatabaseInstanceUpdate = `
{{.BaseNetwork}}
{{.BaseNewFlavor}}

resource "vkcs_db_instance" "basic" {
  name             = "basic"
  flavor_id = data.vkcs_compute_flavor.base.id
  size = 9
  volume_type = "{{.VolumeType}}"

  datastore {
    version = "13"
    type    = "postgresql"
  }

  network {
    uuid = vkcs_networking_network.base.id
  }
  availability_zone = "{{.AvailabilityZone}}"
  floating_ip_enabled = true
  cloud_monitoring_enabled = true

  disk_autoexpand {
    autoexpand = true
    max_disk_size = 2000
  }
  depends_on = [vkcs_networking_router_interface.base]

}
`

const testAccDatabaseInstanceRootUser = `
{{.BaseNetwork}}
{{.BaseFlavor}}

resource "vkcs_db_instance" "basic" {
  name = "basic"
  flavor_id = data.vkcs_compute_flavor.base.id
  size = 10
  volume_type = "{{.VolumeType}}"

  datastore {
    version = "13"
    type    = "postgresql"
  }

  network {
    uuid = vkcs_networking_network.base.id
  }
  root_enabled = true
  depends_on = [vkcs_networking_router_interface.base]
}
`

const testAccDatabaseInstanceWal = `
{{.BaseNetwork}}
{{.BaseFlavor}}

resource "vkcs_db_instance" "basic" {
  name             = "basic_wal"
  flavor_id = data.vkcs_compute_flavor.base.id
  size = 8
  volume_type = "{{.VolumeType}}"

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

  wal_volume {
	  size = 8
	  volume_type = "{{.VolumeType}}"
  }

  wal_disk_autoexpand {
	  autoexpand = true
	  max_disk_size = 1000
  }

  depends_on = [vkcs_networking_router_interface.base]
}
`
