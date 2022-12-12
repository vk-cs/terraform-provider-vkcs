---
layout: "vkcs"
page_title: "vkcs: vkcs_db_cluster"
description: |-
  Manages a db cluster.
---

# vkcs_db_cluster

Provides a db cluster resource. This can be used to create, modify and delete db cluster for galera_mysql, postgresql, tarantool datastores.

## Example Usage
### Basic cluster
```terraform
resource "vkcs_db_cluster" "db-cluster" {
  name        = "db-cluster"

  availability_zone = "GZ1"
  datastore {
    type    = "postgresql"
    version = "12"
  }

  cluster_size = 3

  flavor_id   = data.vkcs_compute_flavor.db.id

  volume_size = 10
  volume_type = "ceph-ssd"

  network {
    uuid = vkcs_networking_network.db.id
  }

  depends_on = [
    vkcs_networking_network.db,
    vkcs_networking_subnet.db
  ]
}
```

### Cluster restored from backup
```terraform
resource "vkcs_db_cluster" "mydb-cluster" {
  name        = "mydb-cluster"

  datastore {
    type    = "postgresql"
    version = "12"
  }

  cluster_size = 3

  flavor_id   = "9e931469-1490-489e-88af-29a289681c53"

  volume_size = 10
  volume_type = "MS1"

  network {
    uuid = "3ee9b184-3311-4d85-840b-7a9c48e7beac"
  }

  restore_point {
    backup_id = "backup_id"
  }
}
```

### Cluster with scheduled PITR backup
```terraform
resource "vkcs_db_cluster" "mydb-cluster" {
  name        = "mydb-cluster"

  datastore {
    type    = "postgresql"
    version = "12"
  }

  cluster_size = 3

  flavor_id   = "9e931469-1490-489e-88af-29a289681c53"

  volume_size = 10
  volume_type = "MS1"

  network {
    uuid = "3ee9b184-3311-4d85-840b-7a9c48e7beac"
  }

  backup_schedule {
    name = three_hours_backup
    start_hours = 16
    start_minutes = 20
    interval_hours = 3
    keep_count = 3
  }
}
```
## Argument Reference
- `cluster_size` **Number** (***Required***) The number of instances in the cluster.

- `datastore` (***Required***) Object that represents datastore of the cluster. Changing this creates a new cluster.
  - `type` **String** (***Required***) Type of the datastore. Changing this creates a new cluster. Type of the datastore can either be "galera_mysql", "postgresql" or "tarantool".

  - `version` **String** (***Required***) Version of the datastore. Changing this creates a new cluster.

- `flavor_id` **String** (***Required***) The ID of flavor for the cluster.

- `name` **String** (***Required***) The name of the cluster. Changing this creates a new cluster.

- `volume_size` **Number** (***Required***) Size of the cluster instance volume.

- `volume_type` **String** (***Required***) The type of the cluster instance volume. Changing this creates a new cluster.

- `availability_zone` **String** (*Optional*) The name of the availability zone of the cluster. Changing this creates a new cluster.

- `backup_schedule` (*Optional*) Object that represents configuration of PITR backup. This functionality is available only for postgres datastore. **New since v.0.1.4**
  - `interval_hours` **Number** (***Required***) Time interval between backups, specified in hours. Available values: 3, 6, 8, 12, 24.

  - `keep_count` **Number** (***Required***) Number of backups to be stored.

  - `name` **String** (***Required***) Name of the schedule.

  - `start_hours` **Number** (***Required***) Hours part of timestamp of initial backup.

  - `start_minutes` **Number** (***Required***) Minutes part of timestamp of initial backup.

- `capabilities` (*Optional*) Object that represents capability applied to cluster. There can be several instances of this object.
  - `name` **String** (***Required***) The name of the capability to apply.

  - `settings` <strong>Map of </strong>**String** (*Optional*) Map of key-value settings of the capability.

- `configuration_id` **String** (*Optional*) The id of the configuration attached to cluster.

- `disk_autoexpand` (*Optional*) Object that represents autoresize properties of the cluster.
  - `autoexpand` **Boolean** (*Optional*) Indicates whether autoresize is enabled.

  - `max_disk_size` **Number** (*Optional*) Maximum disk size for autoresize.

- `floating_ip_enabled` **Boolean** (*Optional*) Indicates whether floating ip is created for cluster. Changing this creates a new cluster.

- `keypair` **String** (*Optional*) Name of the keypair to be attached to cluster. Changing this creates a new cluster.

- `network` (*Optional*) Object that represents network of the cluster. Changing this creates a new cluster.
  - `port` **String** (*Optional*) The port id of the network. Changing this creates a new cluster.

  - `uuid` **String** (*Optional*) The id of the network. Changing this creates a new cluster.

- `region` **String** (*Optional*) Region to create resource in.

- `restore_point` (*Optional*) Object that represents backup to restore cluster from. **New since v.0.1.4**.
  - `backup_id` **String** (***Required***) ID of the backup.

  - `target` **String** (*Optional*) Used only for restoring from PITR backups. Timestamp of needed backup in format "2021-10-06 01:02:00". You can specify "latest" to use most recent backup.

- `root_enabled` **Boolean** (*Optional*) Indicates whether root user is enabled for the cluster.

- `root_password` **String** (*Optional* Sensitive) Password for the root user of the cluster.

- `shrink_options` **String** (*Optional*) Used only for shrinking cluster. List of IDs of instances that should remain after shrink. If no options are supplied, shrink operation will choose first non-leader instance to delete.

- `wal_disk_autoexpand` (*Optional*) Object that represents autoresize properties of wal volume of the cluster.
  - `autoexpand` **Boolean** (*Optional*) Indicates whether wal volume autoresize is enabled.

  - `max_disk_size` **Number** (*Optional*) Maximum disk size for wal volume autoresize.

- `wal_volume` (*Optional*) Object that represents wal volume of the cluster. Changing this creates a new cluster.
  - `size` **Number** (***Required***) Size of the instance wal volume.

  - `volume_type` **String** (***Required***) The type of the cluster wal volume. Changing this creates a new cluster.


## Attributes Reference
- `cluster_size` **Number** See Argument Reference above.

- `datastore`  See Argument Reference above.
  - `type` **String** See Argument Reference above.

  - `version` **String** See Argument Reference above.

- `flavor_id` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `volume_size` **Number** See Argument Reference above.

- `volume_type` **String** See Argument Reference above.

- `availability_zone` **String** See Argument Reference above.

- `backup_schedule`  See Argument Reference above.
  - `interval_hours` **Number** See Argument Reference above.

  - `keep_count` **Number** See Argument Reference above.

  - `name` **String** See Argument Reference above.

  - `start_hours` **Number** See Argument Reference above.

  - `start_minutes` **Number** See Argument Reference above.

- `capabilities`  See Argument Reference above.
  - `name` **String** See Argument Reference above.

  - `settings` <strong>Map of </strong>**String** See Argument Reference above.

- `configuration_id` **String** See Argument Reference above.

- `disk_autoexpand`  See Argument Reference above.
  - `autoexpand` **Boolean** See Argument Reference above.

  - `max_disk_size` **Number** See Argument Reference above.

- `floating_ip_enabled` **Boolean** See Argument Reference above.

- `keypair` **String** See Argument Reference above.

- `network`  See Argument Reference above.
  - `port` **String** See Argument Reference above.

  - `uuid` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `restore_point`  See Argument Reference above.
  - `backup_id` **String** See Argument Reference above.

  - `target` **String** See Argument Reference above.

- `root_enabled` **Boolean** See Argument Reference above.

- `root_password` **String** See Argument Reference above.

- `shrink_options` **String** See Argument Reference above.

- `wal_disk_autoexpand`  See Argument Reference above.
  - `autoexpand` **Boolean** See Argument Reference above.

  - `max_disk_size` **Number** See Argument Reference above.

- `wal_volume`  See Argument Reference above.
  - `size` **Number** See Argument Reference above.

  - `volume_type` **String** See Argument Reference above.

- `id` **String** ID of the resource.

- `instances` **Object** Cluster instances info.



## Import

Clusters can be imported using the `id`, e.g.

```shell
terraform import vkcs_db_cluster.mycluster 708a74a1-6b00-4a96-938c-28a8a6d98590
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.

You should at least add following fields to your .tf file:

`name, flavor_id, cluster_size, volume_size, volume_type, datastore`

Please, use `"IMPORTED"` as value for `volume_type` field.
