---
subcategory: "Databases"
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

  cloud_monitoring_enabled = true

  shard {
    availability_zone = "GZ1"
    size        = 1
    shard_id    = "shard0"
    flavor_id   = data.vkcs_compute_flavor.db.id

    volume_size = 8
    volume_type = "ceph-ssd"

    network {
      uuid = vkcs_networking_network.db.id
      security_groups = [vkcs_networking_secgroup.secgroup.id]
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
      security_groups = [vkcs_networking_secgroup.secgroup.id]
    }
  }

  depends_on = [
    vkcs_networking_router_interface.db,
    vkcs_networking_secgroup.secgroup
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
- `datastore` **required** &rarr;  Object that represents datastore of the cluster. Changing this creates a new cluster.
  - `type` **required** *string* &rarr;  Type of the datastore. Changing this creates a new cluster. Type of the datastore must be "clickhouse".

  - `version` **required** *string* &rarr;  Version of the datastore. Changing this creates a new cluster.

- `name` **required** *string* &rarr;  The name of the cluster. Changing this creates a new cluster.

- `shard` **required** &rarr;  Object that represents cluster shard. There can be several instances of this object.
  - `flavor_id` **required** *string* &rarr;  The ID of flavor for the cluster shard.

  - `shard_id` **required** *string* &rarr;  The ID of the shard. Changing this creates a new cluster.

  - `size` **required** *number* &rarr;  The number of instances in the cluster shard.

  - `volume_size` **required** *number* &rarr;  Size of the cluster shard instance volume.

  - `volume_type` **required** *string* &rarr;  The type of the cluster shard instance volume.

  - `availability_zone` optional *string* &rarr;  The name of the availability zone of the cluster shard. Changing this creates a new cluster.

  - `network` optional
    - `port` optional deprecated *string* &rarr;  The port id of the network. Changing this creates a new cluster. **Deprecated** This argument is deprecated, please do not use it.

    - `security_groups` optional *set of* *string* &rarr;  An array of one or more security group IDs to associate with the shard instances. Changing this creates a new cluster. **New since v0.2.0**.

    - `subnet_id` optional *string* &rarr;  The id of the subnet. Changing this creates a new cluster. **New since v0.1.15**.

    - `uuid` optional *string* &rarr;  The id of the network. Changing this creates a new cluster.**Note** Although this argument is marked as optional, it is actually required at the moment. Not setting a value for it may cause an error.

  - `shrink_options` optional *string* &rarr;  Used only for shrinking cluster. List of IDs of instances that should remain after shrink. If no options are supplied, shrink operation will choose first non-leader instance to delete.

  - `wal_volume` optional &rarr;  Object that represents wal volume of the cluster.
    - `size` **required** *number* &rarr;  Size of the instance wal volume.

    - `volume_type` **required** *string* &rarr;  The type of the cluster wal volume.

- `capabilities` optional &rarr;  Object that represents capability applied to cluster. There can be several instances of this object.
  - `name` **required** *string* &rarr;  The name of the capability to apply.

  - `settings` optional *map of* *string* &rarr;  Map of key-value settings of the capability.

- `cloud_monitoring_enabled` optional *boolean* &rarr;  Enable cloud monitoring for the cluster. Changing this for Redis or MongoDB creates a new instance. **New since v0.2.0**.

- `configuration_id` optional *string* &rarr;  The id of the configuration attached to cluster.

- `disk_autoexpand` optional &rarr;  Object that represents autoresize properties of the cluster.
  - `autoexpand` optional *boolean* &rarr;  Indicates whether autoresize is enabled.

  - `max_disk_size` optional *number* &rarr;  Maximum disk size for autoresize.

- `floating_ip_enabled` optional *boolean* &rarr;  Boolean field that indicates whether floating ip is created for cluster. Changing this creates a new cluster.

- `keypair` optional *string* &rarr;  Name of the keypair to be attached to cluster. Changing this creates a new cluster.

- `region` optional *string* &rarr;  Region to create resource in.

- `restore_point` optional &rarr;  Object that represents backup to restore instance from. **New since v0.1.4**.
  - `backup_id` **required** *string* &rarr;  ID of the backup.

- `root_enabled` optional *boolean* &rarr;  Indicates whether root user is enabled for the cluster.

- `root_password` optional sensitive *string* &rarr;  Password for the root user of the cluster.

- `wal_disk_autoexpand` optional &rarr;  Object that represents autoresize properties of wal volume of the cluster.
  - `autoexpand` optional *boolean* &rarr;  Indicates whether wal volume autoresize is enabled.

  - `max_disk_size` optional *number* &rarr;  Maximum disk size for wal volume autoresize.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `shard` 
  - `instances` *object* &rarr;  Shard instances info. **New since v0.1.15**.



## Import

Clusters can be imported using the `id`, e.g.

```shell
terraform import vkcs_db_cluster_with_shards.mycluster 708a74a1-6b00-4a96-938c-28a8a6d98590
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.

You should at least add following fields to your .tf file:

`name, datastore`, and for each shard add: `shard_id, size, flavor_id, volume_size, volume_type`

Please, use `"IMPORTED"` as value for `volume_type` field.
