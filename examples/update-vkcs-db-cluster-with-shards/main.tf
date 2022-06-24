terraform {
  required_providers {
    vkcs = {
      source  = "vk-cs/vkcs"
      version = "~> 0.1.0"
    }
  }
}

data "vkcs_compute_flavor" "db" {
  name = var.db-instance-flavor
}

resource "vkcs_networking_network" "db" {
  name           = "db-net"
}

resource "vkcs_db_cluster_with_shards" "db-cluster-with-shards" {
  name = "db-cluster-with-shards"

  datastore {
    type    = "clickhouse"
    version = "20"
  }

  shard {
    size        = 2
    shard_id    = "shard0"
    flavor_id   = data.vkcs_compute_flavor.db.id

    volume_size = 12
    volume_type = "ceph-ssd"

    network {
      uuid = vkcs_networking_network.db.id
    }
  }

  shard {
    size        = 2
    shard_id    = "shard1"
    flavor_id   = data.vkcs_compute_flavor.db.id

    volume_size = 12
    volume_type = "ceph-ssd"

    network {
      uuid = vkcs_networking_network.db.id
    }
  }
}