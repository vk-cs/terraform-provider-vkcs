package db_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/instances"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/users"
)

func TestAccDatabaseUser_basic(t *testing.T) {
	var user users.User
	var instance instances.InstanceResp

	baseConfig := acctest.AccTestRenderConfig(testAccDatabaseUserBase)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: baseConfig,
				Check:  testAccCheckDatabaseInstanceExists("vkcs_db_instance.base", &instance),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseUserBasic, map[string]string{"TestAccDatabaseUserBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseUserExists("vkcs_db_user.user", &instance, &user),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "name", "tfacc-basic"),
					resource.TestCheckResourceAttrPair("vkcs_db_user.user", "dbms_id", "vkcs_db_instance.base", "id"),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "dbms_type", "instance"),
				),
			},
			{
				ResourceName:            "vkcs_db_user.user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccDatabaseUser_full(t *testing.T) {
	var user users.User
	var instance instances.InstanceResp

	baseConfig := acctest.AccTestRenderConfig(testAccDatabaseUserBase)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: baseConfig,
				Check:  testAccCheckDatabaseInstanceExists("vkcs_db_instance.base", &instance),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseUserFull, map[string]string{"TestAccDatabaseUserBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseUserExists("vkcs_db_user.user", &instance, &user),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "name", "tfacc-full"),
					resource.TestCheckResourceAttrPair("vkcs_db_user.user", "dbms_id", "vkcs_db_instance.base", "id"),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "host", "192.168.0.1"),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "databases.#", "2"),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "databases.0", "tfacc-db_1"),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "databases.1", "tfacc-db_2"),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "dbms_type", "instance"),
				),
			},
			{
				ResourceName:            "vkcs_db_user.user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"host", "password"},
			},
		},
	})
}

func TestAccDatabaseUser_update(t *testing.T) {
	var user users.User
	var instance instances.InstanceResp

	baseConfig := acctest.AccTestRenderConfig(testAccDatabaseUserBase)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: baseConfig,
				Check:  testAccCheckDatabaseInstanceExists("vkcs_db_instance.base", &instance),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseUpdate, map[string]string{
					"TestAccDatabaseUserBase": baseConfig,
					"Name":                    "tfacc-user",
					"Password":                "Qw!weZ12$234Ax09",
					"Host":                    "192.168.0.1",
					"Databases":               `[ vkcs_db_database.db_1.name, vkcs_db_database.db_2.name ]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseUserExists("vkcs_db_user.user", &instance, &user),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "name", "tfacc-user"),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "host", "192.168.0.1"),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "databases.#", "2"),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "databases.0", "tfacc-db_1"),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "databases.1", "tfacc-db_2"),
				),
			},
			{
				ResourceName:            "vkcs_db_user.user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"host", "password"},
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseUpdate, map[string]string{
					"TestAccDatabaseUserBase": baseConfig,
					"Name":                    "tfacc-new_user",
					"Password":                "rTqn!I24$",
					"Host":                    "192.168.0.2",
					"Databases":               `[ vkcs_db_database.db_3.name, vkcs_db_database.db_2.name ]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseUserExists("vkcs_db_user.user", &instance, &user),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "name", "tfacc-new_user"),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "host", "192.168.0.2"),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "databases.#", "2"),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "databases.0", "tfacc-db_3"),
					resource.TestCheckResourceAttr("vkcs_db_user.user", "databases.1", "tfacc-db_2"),
				),
			},
			{
				ResourceName:            "vkcs_db_user.user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"host", "password"},
			},
		},
	})
}

func TestAccDatabaseUser_skipDeletion(t *testing.T) {
	var user users.User
	var instance instances.InstanceResp

	baseConfig := acctest.AccTestRenderConfig(testAccDatabaseUserBase)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: baseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseInstanceExists("vkcs_db_instance.base", &instance),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseUserSkipDeletion, map[string]string{"TestAccDatabaseUserBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseUserExists("vkcs_db_user.user", &instance, &user),
				),
			},
			{
				ResourceName:            "vkcs_db_user.user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "vendor_options"},
			},
			{
				Config: baseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseUserExistsInAPI(&instance, &user),
				),
			},
		},
	})
}

func testAccCheckDatabaseUserExists(n string, instance *instances.InstanceResp, user *users.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no id is set")
		}

		parts := strings.SplitN(rs.Primary.ID, "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("malformed user name: %s", rs.Primary.ID)
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		DatabaseClient, err := config.DatabaseV1Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating cloud database client: %s", err)
		}

		pages, err := users.List(DatabaseClient, instance.ID, "instance").AllPages()
		if err != nil {
			return fmt.Errorf("unable to retrieve users: %s", err)
		}

		allUsers, err := users.ExtractUsers(pages)
		if err != nil {
			return fmt.Errorf("unable to extract users: %s", err)
		}

		for _, u := range allUsers {
			if u.Name == parts[1] {
				*user = u
				return nil
			}
		}

		return fmt.Errorf("user %s does not exist", n)
	}
}

func testAccCheckDatabaseUserExistsInAPI(instance *instances.InstanceResp, user *users.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.AccTestProvider.Meta().(clients.Config)
		client, err := config.DatabaseV1Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating Databases API client: %s", err)
		}

		pages, err := users.List(client, instance.ID, "instance").AllPages()
		if err != nil {
			return fmt.Errorf("unable to retrieve users: %s", err)
		}

		allUsers, err := users.ExtractUsers(pages)
		if err != nil {
			return fmt.Errorf("unable to extract users: %s", err)
		}

		for _, u := range allUsers {
			if u.Name == user.Name {
				*user = u
				return nil
			}
		}

		return fmt.Errorf("user %s does not exist", user.Name)
	}
}

func testAccCheckDatabaseUserDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)

	DatabaseClient, err := config.DatabaseV1Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("error creating cloud database client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_db_user" {
			continue
		}

		parts := strings.SplitN(rs.Primary.ID, "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("malformed username: %s", rs.Primary.ID)
		}

		pages, err := users.List(DatabaseClient, parts[0], "instance").AllPages()
		if err != nil {
			return nil
		}

		allUsers, err := users.ExtractUsers(pages)
		if err != nil {
			return fmt.Errorf("unable to extract users: %s", err)
		}

		var exists bool
		for _, v := range allUsers {
			if v.Name == parts[1] {
				exists = true
			}
		}

		if exists {
			return fmt.Errorf("user still exists")
		}
	}

	return nil
}

const testAccDatabaseUserBase = `
{{ .BaseNetwork }}
{{ .BaseFlavor }}

resource "vkcs_db_instance" "base" {
  availability_zone = "{{ .AvailabilityZone }}"

  name        = "tfacc-base"
  flavor_id   = data.vkcs_compute_flavor.base.id
  size        = 10
  volume_type = "{{.VolumeType}}"

  datastore {
    version = "13"
    type    = "postgresql"
  }

  network {
    uuid = vkcs_networking_network.base.id
  }

  depends_on = [
    vkcs_networking_router_interface.base
  ]
}

resource "vkcs_db_database" "db_1" {
  name    = "tfacc-db_1"
  dbms_id = vkcs_db_instance.base.id
}

resource "vkcs_db_database" "db_2" {
  name    = "tfacc-db_2"
  dbms_id = vkcs_db_instance.base.id
}
`

const testAccDatabaseUserBasic = `
{{ .TestAccDatabaseUserBase }}

resource "vkcs_db_user" "user" {
  name        = "tfacc-basic"
  dbms_id     = vkcs_db_instance.base.id
  password    = "Qw!weZ12$234Ax09"
}
`

const testAccDatabaseUserFull = `
{{ .TestAccDatabaseUserBase }}

resource "vkcs_db_user" "user" {
  name     = "tfacc-full"
  dbms_id  = vkcs_db_instance.base.id
  password = "Qw!weZ12$234Ax09"
  host     = "192.168.0.1"
  databases = [
    vkcs_db_database.db_1.name,
    vkcs_db_database.db_2.name
  ]
}
`

const testAccDatabaseUpdate = `
{{ .TestAccDatabaseUserBase }}

resource "vkcs_db_database" "db_3" {
  name    = "tfacc-db_3"
  dbms_id = vkcs_db_instance.base.id
}

resource "vkcs_db_user" "user" {
  name      = "{{ .Name }}"
  dbms_id   = vkcs_db_instance.base.id
  password  = "{{ .Password }}"
  host      = "{{ .Host }}"
  databases = {{ .Databases }}
}
`

const testAccDatabaseUserSkipDeletion = `
{{ .TestAccDatabaseUserBase }}

resource "vkcs_db_user" "user" {
  name     = "tfacc-basic"
  dbms_id  = vkcs_db_instance.base.id
  password = "Qw!weZ12$234Ax09"

  vendor_options {
    skip_deletion = true
  }
}
`
