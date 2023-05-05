package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/clusters"
)

func TestAccDatabaseClusterWithShards_basic(t *testing.T) {
	var cluster clusters.ClusterResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseClusterWithShardsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccDatabaseClusterWithShardsBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists("vkcs_db_cluster_with_shards.basic", &cluster),
					resource.TestCheckResourceAttrPtr("vkcs_db_cluster_with_shards.basic", "name", &cluster.Name),
				),
			},
		},
	})
}

func TestAccDatabaseClusterWithShards_update(t *testing.T) {
	var cluster clusters.ClusterResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseClusterWithShardsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccDatabaseClusterWithShardsUpdateInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists("vkcs_db_cluster_with_shards.update", &cluster),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.update", "name", "update"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.update", "datastore.0.version", "20.8"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.update", "datastore.0.type", "clickhouse"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.update", "cloud_monitoring_enabled", "false"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.update", "shard.#", "1"),
					resource.TestCheckResourceAttrPair("vkcs_db_cluster_with_shards.update", "shard.0.flavor_id",
						"data.vkcs_compute_flavor.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_db_cluster_with_shards.update", "shard.0.network.0.uuid",
						"vkcs_networking_network.base", "id"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.update", "shard.0.volume_size", "8"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.update", "shard.0.volume_type", "ceph-ssd"),
				),
			},
			{
				Config: testAccRenderConfig(testAccDatabaseClusterWithShardsUpdateUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.update", "name", "update"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.update", "datastore.0.version", "20.8"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.update", "datastore.0.type", "clickhouse"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.update", "cloud_monitoring_enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.update", "shard.#", "1"),
					resource.TestCheckResourceAttrPair("vkcs_db_cluster_with_shards.update", "shard.0.flavor_id",
						"data.vkcs_compute_flavor.new_flavor", "id"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.update", "shard.0.volume_size", "10"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.update", "shard.0.volume_type", "ceph-hdd"),
				),
			},
		},
	})
}

func TestAccDatabaseClusterWithShards_resize(t *testing.T) {
	var cluster clusters.ClusterResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseClusterWithShardsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccDatabaseClusterWithShardsResizeInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists("vkcs_db_cluster_with_shards.resize", &cluster),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.#", "1"),

					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.0.size", "1"),
					resource.TestCheckResourceAttrPair("vkcs_db_cluster_with_shards.resize", "shard.0.flavor_id",
						"data.vkcs_compute_flavor.base", "id"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.0.volume_size", "8"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.0.volume_type", "ceph-ssd"),
					resource.TestCheckResourceAttrPair("vkcs_db_cluster_with_shards.resize", "shard.0.network.0.uuid",
						"vkcs_networking_network.base", "id"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.0.instances.#", "1"),
				),
			},
			{
				Config: testAccRenderConfig(testAccDatabaseClusterWithShardsResizeGrow),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists("vkcs_db_cluster_with_shards.resize", &cluster),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.#", "2"),

					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.0.size", "3"),
					resource.TestCheckResourceAttrPair("vkcs_db_cluster_with_shards.resize", "shard.0.flavor_id",
						"data.vkcs_compute_flavor.new_flavor", "id"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.0.volume_size", "10"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.0.volume_type", "ceph-hdd"),
					resource.TestCheckResourceAttrPair("vkcs_db_cluster_with_shards.resize", "shard.0.network.0.uuid",
						"vkcs_networking_network.base", "id"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.0.instances.#", "3"),

					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.1.size", "1"),
					resource.TestCheckResourceAttrPair("vkcs_db_cluster_with_shards.resize", "shard.1.flavor_id",
						"data.vkcs_compute_flavor.base", "id"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.1.volume_size", "8"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.1.volume_type", "ceph-ssd"),
					resource.TestCheckResourceAttrPair("vkcs_db_cluster_with_shards.resize", "shard.1.network.0.uuid",
						"vkcs_networking_network.base", "id"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.1.instances.#", "1"),
				),
			},
			{
				Config: testAccRenderConfig(testAccDatabaseClusterWithShardsResizeShrink),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists("vkcs_db_cluster_with_shards.resize", &cluster),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.#", "1"),

					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.0.size", "1"),
					resource.TestCheckResourceAttrPair("vkcs_db_cluster_with_shards.resize", "shard.0.flavor_id",
						"data.vkcs_compute_flavor.base", "id"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.0.volume_size", "10"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.0.volume_type", "ceph-ssd"),
					resource.TestCheckResourceAttrPair("vkcs_db_cluster_with_shards.resize", "shard.0.network.0.uuid",
						"vkcs_networking_network.base", "id"),
					resource.TestCheckResourceAttr("vkcs_db_cluster_with_shards.resize", "shard.0.instances.#", "1"),
				),
			},
		},
	})
}

func testAccCheckDatabaseClusterWithShardsDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(clients.Config)

	DatabaseClient, err := config.DatabaseV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS database client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_db_cluster_with_shards" {
			continue
		}
		_, err := clusters.Get(DatabaseClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("cluster still exists")
		}
	}

	return nil
}

const testAccDatabaseClusterWithShardsBasic = `
{{.BaseNetwork}}
{{.BaseFlavor}}

 resource "vkcs_db_cluster_with_shards" "basic" {
	name      = "basic"

	datastore {
	  version = "20.8"
	  type    = "clickhouse"
	}
  
	shard {
	  size = 1
	  shard_id = "shard0"
	  flavor_id = data.vkcs_compute_flavor.base.id
	  volume_size      = 8
	  volume_type = "{{.VolumeType}}"
	  network {
		  uuid = vkcs_networking_network.base.id
	  }
	  availability_zone = "{{.AvailabilityZone}}"
	}

	depends_on = [vkcs_networking_router_interface.base]
 }
`

const testAccDatabaseClusterWithShardsUpdateInitial = `
{{.BaseNetwork}}
{{.BaseFlavor}}

resource "vkcs_db_cluster_with_shards" "update" {
  name      = "update"
  
  datastore {
	version = "20.8"
	type    = "clickhouse"
  }

  cloud_monitoring_enabled = false
  
  shard {
	size = 1
	shard_id = "shard0"
	flavor_id = data.vkcs_compute_flavor.base.id
	volume_size = 8
	volume_type = "ceph-ssd"
	network {
	  uuid = vkcs_networking_network.base.id
	}
	availability_zone = "{{.AvailabilityZone}}"
  }
  
  depends_on = [vkcs_networking_router_interface.base]
}
`

const testAccDatabaseClusterWithShardsUpdateUpdated = `
{{.BaseNetwork}}
{{.BaseFlavor}}

data "vkcs_compute_flavor" "new_flavor" {
	name = "Standard-4-8-80"
}

resource "vkcs_db_cluster_with_shards" "update" {
  name      = "update"
  
  datastore {
	version = "20.8"
	type    = "clickhouse"
  }

  disk_autoexpand {
	autoexpand = true
	max_disk_size = 1000
  }

  cloud_monitoring_enabled = true
  
  shard {
	size = 1
	shard_id = "shard0"
	flavor_id = data.vkcs_compute_flavor.new_flavor.id
	volume_size = 10
	volume_type = "ceph-hdd"
	network {
	  uuid = vkcs_networking_network.base.id
	}
	availability_zone = "{{.AvailabilityZone}}"
  }

  depends_on = [vkcs_networking_router_interface.base]
}
`

const testAccDatabaseClusterWithShardsResizeInitial = `
{{.BaseNetwork}}
{{.BaseFlavor}}

resource "vkcs_db_cluster_with_shards" "resize" {
  name      = "resize"
  
  datastore {
	version = "20.8"
	type    = "clickhouse"
  }
  
  shard {
	size = 1
	shard_id = "shard0"
	flavor_id = data.vkcs_compute_flavor.base.id
	volume_size = 8
	volume_type = "ceph-ssd"
	network {
	  uuid = vkcs_networking_network.base.id
	}
	availability_zone = "{{.AvailabilityZone}}"
  }
  
  depends_on = [vkcs_networking_router_interface.base]
}
`

const testAccDatabaseClusterWithShardsResizeGrow = `
{{.BaseNetwork}}
{{.BaseFlavor}}

data "vkcs_compute_flavor" "new_flavor" {
	name = "Standard-4-8-80"
}

resource "vkcs_db_cluster_with_shards" "resize" {
  name      = "resize"
  
  datastore {
	version = "20.8"
	type    = "clickhouse"
  }
  
  shard {
	size = 3
	shard_id = "shard0"
	flavor_id = data.vkcs_compute_flavor.new_flavor.id
	volume_size = 10
	volume_type = "ceph-hdd"
	network {
	  uuid = vkcs_networking_network.base.id
	}
	availability_zone = "{{.AvailabilityZone}}"
  }

  shard {
	size = 1
	shard_id = "shard1"
	flavor_id = data.vkcs_compute_flavor.base.id
	volume_size = 8
	volume_type = "ceph-ssd"
	network {
	  uuid = vkcs_networking_network.base.id
	}
	availability_zone = "{{.AvailabilityZone}}"
  }
  
  depends_on = [vkcs_networking_router_interface.base]
}
`

const testAccDatabaseClusterWithShardsResizeShrink = `
{{.BaseNetwork}}
{{.BaseFlavor}}

resource "vkcs_db_cluster_with_shards" "resize" {
  name      = "resize"
  
  datastore {
	version = "20.8"
	type    = "clickhouse"
  }
  
  shard {
	size = 1
	shard_id = "shard0"
	flavor_id = data.vkcs_compute_flavor.base.id
	volume_size = 10
	volume_type = "ceph-ssd"
	network {
	  uuid = vkcs_networking_network.base.id
	}
	availability_zone = "{{.AvailabilityZone}}"
  }
  
  depends_on = [vkcs_networking_router_interface.base]
}
`
