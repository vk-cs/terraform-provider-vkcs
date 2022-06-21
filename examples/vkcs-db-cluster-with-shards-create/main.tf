resource "vkcs_db_cluster_with_shards" "db-cluster-with-shards" {
  name = "db-cluster-with-shards"



  datastore {
    type    = "clickhouse"
    version = "20.8"
  }

  shard {
    availability_zone = "GZ1"
    size        = 1
    shard_id    = "shard0"
    flavor_id   = data.vkcs_compute_flavor.db.id

    volume_size = 8
    volume_type = "ceph-ssd"

    network {
      uuid = vkcs_networking_network.db.id
    }
  }

  shard {
    availability_zone = "GZ1"
    size        = 1
    shard_id    = "shard1"
    flavor_id   = data.vkcs_compute_flavor.db.id
    
    volume_size = 8
    volume_type = "ceph-ssd"

    network {
      uuid = vkcs_networking_network.db.id
    }
  }
}