package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDatabaseCluster_basic(t *testing.T) {
	var cluster dbClusterResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccDatabaseClusterBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists(
						"vkcs_db_cluster.basic", &cluster),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_cluster.basic", "name", &cluster.Name),
				),
			},
			{
				Config: testAccRenderConfig(testAccDatabaseClusterUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists(
						"vkcs_db_cluster.basic", &cluster),
					resource.TestCheckResourceAttr(
						"vkcs_db_cluster.basic", "volume_size", "9"),
				),
			},
		},
	})
}

func TestAccDatabaseCluster_wal(t *testing.T) {
	var cluster dbClusterResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccDatabaseClusterWal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists(
						"vkcs_db_cluster.basic", &cluster),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_cluster.basic", "name", &cluster.Name),
				),
			},
		},
	})
}

func TestAccDatabaseCluster_wal_no_update(t *testing.T) {
	var cluster dbClusterResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccDatabaseClusterWal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists(
						"vkcs_db_cluster.basic", &cluster),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_cluster.basic", "name", &cluster.Name),
				),
			},
			{
				Config: testAccRenderConfig(testAccDatabaseClusterWal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists(
						"vkcs_db_cluster.basic", &cluster),
				),
			},
		},
	})
}

func testAccCheckDatabaseClusterExists(n string, cluster *dbClusterResp) resource.TestCheckFunc {
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

		found, err := dbClusterGet(DatabaseClient, rs.Primary.ID).extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("cluster not found")
		}

		*cluster = *found

		return nil
	}
}

func testAccCheckDatabaseClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)

	DatabaseClient, err := config.DatabaseV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS database client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_db_cluster" {
			continue
		}
		_, err := dbClusterGet(DatabaseClient, rs.Primary.ID).extract()
		if err == nil {
			return fmt.Errorf("cluster still exists")
		}
	}

	return nil
}

const testAccDatabaseClusterBasic = `
{{.BaseNetwork}}
{{.BaseFlavor}}

 resource "vkcs_db_cluster" "basic" {
   name      = "basic"
   flavor_id = data.vkcs_compute_flavor.base.id
   volume_size      = 8
   volume_type = "{{.VolumeType}}"
   cluster_size = 3
   datastore {
	version = "13"
	type    = "postgresql"
  }

   network {
     uuid = vkcs_networking_network.base.id
   }
	
   availability_zone = "{{.AvailabilityZone}}"

   depends_on = [
    vkcs_networking_network.base,
    vkcs_networking_subnet.base
  ]
 }
`

const testAccDatabaseClusterUpdate = `
{{.BaseNetwork}}
{{.BaseNewFlavor}}

resource "vkcs_db_cluster" "basic" {
	name      = "basic"
	flavor_id = data.vkcs_compute_flavor.base.id
	volume_size      = 9
	volume_type = "{{.VolumeType}}"
	cluster_size = 3
	datastore {
	 version = "13"
	 type    = "postgresql"
   }
 
	network {
	  uuid = vkcs_networking_network.base.id
	}
	 
	availability_zone = "{{.AvailabilityZone}}"

	depends_on = [
		vkcs_networking_network.base,
		vkcs_networking_subnet.base
	  ]
  }
`

const testAccDatabaseClusterWal = `
{{.BaseNetwork}}
				
{{.BaseFlavor}}

 resource "vkcs_db_cluster" "basic" {
   name      = "basic"
   flavor_id = data.vkcs_compute_flavor.base.id
   volume_size      = 8
   volume_type = "{{.VolumeType}}"
   cluster_size = 3
   datastore {
	version = "13"
	type    = "postgresql"
  }

   network {
     uuid = vkcs_networking_network.base.id
   }
	
   availability_zone = "{{.AvailabilityZone}}"
   wal_volume {
	size = 8
	volume_type = "{{.VolumeType}}"
   }

   wal_disk_autoexpand {
	autoexpand = true
	max_disk_size = 1000
   }

   depends_on = [
    vkcs_networking_network.base,
    vkcs_networking_subnet.base
  ]
 }
`
