---
layout: "vkcs"
page_title: "vkcs: vkcs_db_cluster_with_shards"
description: |-
  Manages a db cluster with shards.
---

# vkcs_db_cluster_with_shards

Provides a db cluster with shards resource. This can be used to create, modify and delete db cluster with shards for clickhouse datastore.

## Example Usage
### Basic cluster with shards
```terraform
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

  depends_on = [
    vkcs_networking_router_interface.db
  ]
}

locals {
  cluster = vkcs_db_cluster_with_shards.db-cluster-with-shards
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
```

### Cluster with shards restored from backup
```terraform
resource "vkcs_db_cluster_with_shards" "db-cluster-with-shards" {
  name = "db-cluster-with-shards"

  datastore {
    type    = "clickhouse"
    version = "20.8"
  }

  shard {
    size        = 2
    shard_id    = "shard0"
    flavor_id   = "9e931469-1490-489e-88af-29a289681c53"

    volume_size = 10
    volume_type = "ceph-ssd"
    
    network {
      uuid = "3ee9b184-3311-4d85-840b-7a9c48e7beac"
    }
  }

  shard {
    size        = 2
    shard_id    = "shard1"
    flavor_id   = "9e931469-1490-489e-88af-29a289681c53"
    
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
```
## Argument Reference
- `datastore` (***Required***) Object that represents datastore of the cluster. Changing this creates a new cluster.
  - `type` **String** (***Required***) Type of the datastore. Changing this creates a new cluster. Type of the datastore must be "clickhouse".

  - `version` **String** (***Required***) Version of the datastore. Changing this creates a new cluster.

- `name` **String** (***Required***) The name of the cluster. Changing this creates a new cluster.

- `shard` (***Required***) Object that represents cluster shard. There can be several instances of this object.
  - `flavor_id` **String** (***Required***) The ID of flavor for the cluster shard.

  - `shard_id` **String** (***Required***) The ID of the shard. Changing this creates a new cluster.

  - `size` **Number** (***Required***) The number of instances in the cluster shard.

  - `volume_size` **Number** (***Required***) Size of the cluster shard instance volume.

  - `volume_type` **String** (***Required***) The type of the cluster shard instance volume.

  - `availability_zone` **String** (*Optional*) The name of the availability zone of the cluster shard. Changing this creates a new cluster.

  - `network` (*Optional*)
    - `port` **String** (*Optional* Deprecated) The port id of the network. Changing this creates a new cluster. ***Deprecated*** This argument is deprecated, please do not use it.

    - `subnet_id` **String** (*Optional*) The id of the subnet. Changing this creates a new cluster. **New since v.0.1.15**.

    - `uuid` **String** (*Optional*) The id of the network. Changing this creates a new cluster.**Note** Although this argument is marked as optional, it is actually required at the moment. Not setting a value for it may cause an error.

  - `wal_volume` (*Optional*) Object that represents wal volume of the cluster.
    - `size` **Number** (***Required***) Size of the instance wal volume.

    - `volume_type` **String** (***Required***) The type of the cluster wal volume.

- `capabilities` (*Optional*) Object that represents capability applied to cluster. There can be several instances of this object.
  - `name` **String** (***Required***) The name of the capability to apply.

  - `settings` <strong>Map of </strong>**String** (*Optional*) Map of key-value settings of the capability.

- `configuration_id` **String** (*Optional*) The id of the configuration attached to cluster.

- `disk_autoexpand` (*Optional*) Object that represents autoresize properties of the cluster.
  - `autoexpand` **Boolean** (*Optional*) Indicates whether autoresize is enabled.

  - `max_disk_size` **Number** (*Optional*) Maximum disk size for autoresize.

- `floating_ip_enabled` **Boolean** (*Optional*) Boolean field that indicates whether floating ip is created for cluster. Changing this creates a new cluster.

- `keypair` **String** (*Optional*) Name of the keypair to be attached to cluster. Changing this creates a new cluster.

- `region` **String** (*Optional*) Region to create resource in.

- `restore_point` (*Optional*) Object that represents backup to restore instance from. **New since v.0.1.4**.
  - `backup_id` **String** (***Required***) ID of the backup.

- `root_enabled` **Boolean** (*Optional*) Indicates whether root user is enabled for the cluster.

- `root_password` **String** (*Optional* Sensitive) Password for the root user of the cluster.

- `wal_disk_autoexpand` (*Optional*) Object that represents autoresize properties of wal volume of the cluster.
  - `autoexpand` **Boolean** (*Optional*) Indicates whether wal volume autoresize is enabled.

  - `max_disk_size` **Number** (*Optional*) Maximum disk size for wal volume autoresize.


## Attributes Reference
- `datastore`  See Argument Reference above.
  - `type` **String** See Argument Reference above.

  - `version` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `shard`  See Argument Reference above.
  - `flavor_id` **String** See Argument Reference above.

  - `shard_id` **String** See Argument Reference above.

  - `size` **Number** See Argument Reference above.

  - `volume_size` **Number** See Argument Reference above.

  - `volume_type` **String** See Argument Reference above.

  - `availability_zone` **String** See Argument Reference above.

  - `network` 
    - `port` **String** See Argument Reference above.

    - `subnet_id` **String** See Argument Reference above.

    - `uuid` **String** See Argument Reference above.

  - `wal_volume`  See Argument Reference above.
    - `size` **Number** See Argument Reference above.

    - `volume_type` **String** See Argument Reference above.

  - `instances` **Object** Shard instances info. **New since v.0.1.15**.

- `capabilities`  See Argument Reference above.
  - `name` **String** See Argument Reference above.

  - `settings` <strong>Map of </strong>**String** See Argument Reference above.

- `configuration_id` **String** See Argument Reference above.

- `disk_autoexpand`  See Argument Reference above.
  - `autoexpand` **Boolean** See Argument Reference above.

  - `max_disk_size` **Number** See Argument Reference above.

- `floating_ip_enabled` **Boolean** See Argument Reference above.

- `keypair` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `restore_point`  See Argument Reference above.
  - `backup_id` **String** See Argument Reference above.

- `root_enabled` **Boolean** See Argument Reference above.

- `root_password` **String** See Argument Reference above.

- `wal_disk_autoexpand`  See Argument Reference above.
  - `autoexpand` **Boolean** See Argument Reference above.

  - `max_disk_size` **Number** See Argument Reference above.

- `id` **String** ID of the resource.



## Import

Clusters can be imported using the `id`, e.g.

```shell
terraform import vkcs_db_cluster_with_shards.mycluster 708a74a1-6b00-4a96-938c-28a8a6d98590
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.

You should at least add following fields to your .tf file:

`name, datastore`, and for each shard add: `shard_id, size, flavor_id, volume_size, volume_type`

Please, use `"IMPORTED"` as value for `volume_type` field.
