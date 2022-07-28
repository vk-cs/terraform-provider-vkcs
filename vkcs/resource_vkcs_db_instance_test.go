package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDatabaseInstance_basic(t *testing.T) {
	var instance instanceResp

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseInstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseInstanceExists(
						"vkcs_db_instance.basic", &instance),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_instance.basic", "name", &instance.Name),
				),
			},
			{
				Config: testAccDatabaseInstanceUpdate,
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
	var instance instanceResp
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseInstanceRootUser,
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
	var instance instanceResp

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseInstanceWal,
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
	var instance instanceResp

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseInstanceWal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseInstanceExists(
						"vkcs_db_instance.basic", &instance),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_instance.basic", "name", &instance.Name),
				),
			},
			{
				Config: testAccDatabaseInstanceWal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseInstanceExists(
						"vkcs_db_instance.basic", &instance),
				),
			},
		},
	})
}

func testAccCheckDatabaseInstanceExists(n string, instance *instanceResp) resource.TestCheckFunc {
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

		found, err := instanceGet(DatabaseClient, rs.Primary.ID).extract()
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
	config := testAccProvider.Meta().(configer)

	DatabaseClient, err := config.DatabaseV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS database client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_db_instance" {
			continue
		}
		_, err := instanceGet(DatabaseClient, rs.Primary.ID).extract()
		if err == nil {
			return fmt.Errorf("instance still exists")
		}
	}

	return nil
}

func testAccCheckDatabaseRootUserExists(n string, instance *instanceResp) resource.TestCheckFunc {

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
			return fmt.Errorf("error creating cloud database client: %s", err)
		}

		isRootEnabledResult := instanceRootUserGet(DatabaseClient, rs.Primary.ID)
		isRootEnabled, err := isRootEnabledResult.extract()
		if err != nil {
			return fmt.Errorf("error checking if root user is enabled for instance: %s: %s", rs.Primary.ID, err)
		}

		if isRootEnabled {
			return nil
		}

		return fmt.Errorf("root user %s does not exist", n)
	}
}

var testAccDatabaseInstanceBasic = fmt.Sprintf(`
%s

%s

resource "vkcs_db_instance" "basic" {
  name             = "basic"
  flavor_id = data.vkcs_compute_flavor.base.id
  size = 8
  volume_type = "ceph-ssd"

  datastore {
    version = "13"
    type    = "postgresql"
  }

  network {
    uuid = vkcs_networking_network.base.id
  }
  availability_zone = "GZ1"
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
`, testAccBaseFlavor, testAccBaseNetwork)

var testAccDatabaseInstanceUpdate = fmt.Sprintf(`
%s

%s

resource "vkcs_db_instance" "basic" {
  name             = "basic"
  flavor_id = data.vkcs_compute_flavor.base.id
  size = 9
  volume_type = "ceph-ssd"

  datastore {
    version = "13"
    type    = "postgresql"
  }

  network {
    uuid = vkcs_networking_network.base.id
  }
  availability_zone = "GZ1"
  floating_ip_enabled = true

  disk_autoexpand {
    autoexpand = true
    max_disk_size = 2000
  }
  depends_on = [
    vkcs_networking_network.base,
    vkcs_networking_subnet.base
  ]

}
`, testAccBaseFlavorSecond, testAccBaseNetwork)

var testAccDatabaseInstanceRootUser = fmt.Sprintf(`
%s

%s

resource "vkcs_db_instance" "basic" {
  name = "basic"
  flavor_id = data.vkcs_compute_flavor.base.id
  size = 10
  volume_type = "ceph-ssd"

  datastore {
    version = "13"
    type    = "postgresql"
  }

  network {
    uuid = vkcs_networking_network.base.id
  }
  root_enabled = true
  depends_on = [
    vkcs_networking_network.base,
    vkcs_networking_subnet.base
  ]
}
`, testAccBaseFlavor, testAccBaseNetwork)

var testAccDatabaseInstanceWal = fmt.Sprintf(`
%s

%s

resource "vkcs_db_instance" "basic" {
  name             = "basic_wal"
  flavor_id = data.vkcs_compute_flavor.base.id
  size = 8
  volume_type = "ceph-ssd"

  datastore {
    version = "13"
    type    = "postgresql"
  }

  network {
    uuid = vkcs_networking_network.base.id
  }
  availability_zone = "GZ1"
  floating_ip_enabled = true

  disk_autoexpand {
    autoexpand = true
    max_disk_size = 1000
  }

  wal_volume {
	  size = 8
	  volume_type = "ceph-ssd"
  }

  wal_disk_autoexpand {
	  autoexpand = true
	  max_disk_size = 1000
  }

}
`, testAccBaseFlavor, testAccBaseNetwork)
