package db_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/db"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/clusters"
)

func TestAccDatabaseCluster_basic_big(t *testing.T) {
	var cluster clusters.ClusterResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseClusterBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists(
						"vkcs_db_cluster.basic", &cluster),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_cluster.basic", "name", &cluster.Name),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseClusterUpdate),
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

func TestAccDatabaseCluster_wal_big(t *testing.T) {
	var cluster clusters.ClusterResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseClusterWal),
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

func TestAccDatabaseCluster_wal_no_update_big(t *testing.T) {
	var cluster clusters.ClusterResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseClusterWal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists(
						"vkcs_db_cluster.basic", &cluster),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_cluster.basic", "name", &cluster.Name),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseClusterWal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists(
						"vkcs_db_cluster.basic", &cluster),
				),
			},
		},
	})
}

func TestAccDatabaseCluster_shrink_big(t *testing.T) {
	var cluster clusters.ClusterResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseClusterShrinkInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists(
						"vkcs_db_cluster.basic", &cluster),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseClusterShrinkUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists(
						"vkcs_db_cluster.basic", &cluster),
					testAccCheckDatabaseClusterLeaderExists(&cluster),
				),
			},
		},
	})
}

func testAccCheckDatabaseClusterExists(n string, cluster *clusters.ClusterResp) resource.TestCheckFunc {
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

		found, err := clusters.Get(DatabaseClient, rs.Primary.ID).Extract()
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
	config := acctest.AccTestProvider.Meta().(clients.Config)

	DatabaseClient, err := config.DatabaseV1Client(acctest.OsRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS database client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_db_cluster" {
			continue
		}
		_, err := clusters.Get(DatabaseClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("cluster still exists")
		}
	}

	return nil
}

func testAccCheckDatabaseClusterLeaderExists(cluster *clusters.ClusterResp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, inst := range cluster.Instances {
			if inst.Role == db.DBClusterInstanceRoleLeader {
				return nil
			}
		}
		return fmt.Errorf("cluster leader instance is absent")
	}
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
    vkcs_networking_router_interface.base
  ]
 }
`

const testAccDatabaseClusterUpdate = `
{{.BaseNetwork}}
{{.BaseNewFlavor}}

resource "vkcs_db_cluster" "basic" {
	name      = "basic"
	flavor_id = data.vkcs_compute_flavor.base.id
	cloud_monitoring_enabled = true
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

	depends_on = [vkcs_networking_router_interface.base]
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

   depends_on = [vkcs_networking_router_interface.base]
 }
`

const testAccDatabaseClusterShrinkInitial = `
{{.BaseNetwork}}
				
{{.BaseFlavor}}

 resource "vkcs_db_cluster" "basic" {
   name      = "basic"
   flavor_id = data.vkcs_compute_flavor.base.id
   volume_size      = 8
   volume_type = "ceph-ssd"
   cluster_size = 4
   datastore {
	version = "13"
	type    = "postgresql"
  }
   network {
     uuid = vkcs_networking_network.base.id
   }
	
   availability_zone = "GZ1"
   depends_on = [vkcs_networking_router_interface.base]
 }
`

var testAccDatabaseClusterShrinkUpdated = `
{{.BaseNetwork}}
				
{{.BaseFlavor}}

 resource "vkcs_db_cluster" "basic" {
   name      = "basic"
   flavor_id = data.vkcs_compute_flavor.base.id
   volume_size      = 8
   volume_type = "ceph-ssd"
   cluster_size = 3
   datastore {
	version = "13"
	type    = "postgresql"
  }
   network {
     uuid = vkcs_networking_network.base.id
   }
	
   availability_zone = "GZ1"
   depends_on = [vkcs_networking_router_interface.base]
 }
`
