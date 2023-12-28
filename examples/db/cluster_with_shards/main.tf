resource "vkcs_db_cluster_with_shards" "clickhouse-cluster" {
  name = "clickhouse-cluster-with-shards"

  datastore {
    type    = "clickhouse"
    version = "20.8"
  }

  cloud_monitoring_enabled = true

  shard {
    availability_zone = "GZ1"
    size        = 1
    shard_id    = "shard0"
    flavor_id   = data.vkcs_compute_flavor.basic.id

    volume_size = 8
    volume_type = "ceph-ssd"

    network {
      uuid = vkcs_networking_network.db.id
      security_groups = [vkcs_networking_secgroup.admin.id]
    }
  }

  shard {
    availability_zone = "GZ1"
    size        = 1
    shard_id    = "shard1"
    flavor_id   = data.vkcs_compute_flavor.basic.id

    volume_size = 8
    volume_type = "ceph-ssd"

    network {
      uuid = vkcs_networking_network.db.id
      security_groups = [vkcs_networking_secgroup.admin.id]
    }
  }

  depends_on = [
    vkcs_networking_router_interface.db,
    vkcs_networking_secgroup.admin
  ]
}

locals {
  cluster = vkcs_db_cluster_with_shards.clickhouse-cluster
  shards_ips = {
    for shard in local.cluster.shard : shard.shard_id => [for i in shard.instances : {
      "internal_ip" = i.ip[0]
      "external_ip" = length(i.ip) > 1 ? i.ip[1] : null
    }]
  }
}

output "shard0-ips" {
  value = local.shards_ips["shard0"]
  description = "IPs of instances in shard with \"id\" = \"shard0\""
}
