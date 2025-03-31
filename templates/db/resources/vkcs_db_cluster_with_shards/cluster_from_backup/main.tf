resource "vkcs_db_cluster_with_shards" "db_cluster_with_shards" {
  name = "db-cluster-with-shards"

  datastore {
    type    = "clickhouse"
    version = "24.3"
  }

  shard {
    size      = 2
    shard_id  = "shard0"
    flavor_id = "9e931469-1490-489e-88af-29a289681c53"

    volume_size = 10
    volume_type = "ceph-ssd"

    network {
      uuid = "3ee9b184-3311-4d85-840b-7a9c48e7beac"
    }
  }

  shard {
    size      = 2
    shard_id  = "shard1"
    flavor_id = "9e931469-1490-489e-88af-29a289681c53"

    volume_size = 10
    volume_type = "ceph-ssd"

    network {
      uuid = "3ee9b184-3311-4d85-840b-7a9c48e7beac"
    }

    restore_point {
      backup_id = "7c8110f3-6f7f-4dc3-85c2-16feef9ddc2b"
    }
  }
}
