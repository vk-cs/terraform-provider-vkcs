---
layout: "vkcs"
page_title: "vkcs: db_cluster"
subcategory: ""
description: |-
  Manages a db cluster.
---

# vkcs\_db\_cluster (Resource)

Provides a db cluster resource. This can be used to create, modify and delete db cluster for galera_mysql, postgresql, tarantool datastores.

## Example Usage
### Basic cluster
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
  volume_type = "ceph-ssd"

  network {
    uuid = "3ee9b184-3311-4d85-840b-7a9c48e7beac"
  }
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
  volume_type = "ceph-ssd"

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
  volume_type = "ceph-ssd"

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

The following arguments are supported:

* `name` - (Required) The name of the cluster. Changing this creates a new cluster

* `datastore` - (Required) Object that represents datastore of the cluster. Changing this creates a new cluster. It has following attributes:
    * `type` - (Required) Type of the datastore. Changing this creates a new cluster. Type of the datastore can either be "galera_mysql", "postgresql" or "tarantool".
    * `version` - (Required) Version of the datastore. Changing this creates a new cluster.

* `cluster_size` - (Required) The number of instances in the cluster.

* `keypair` - Name of the keypair to be attached to cluster. Changing this creates a new cluster.

* `floating_ip_enabled` - Boolean field that indicates whether floating ip is created for cluster. Changing this creates a new cluster.

* `flavor_id` - (Required) The ID of flavor for the cluster.

* `availability_zone` - The name of the availability zone of the cluster. Changing this creates a new cluster.

* `volume_size` - (Required) Size of the cluster instance volume.

* `volume_type` - (Required) The type of the cluster instance volume. Changing this creates a new cluster.

* `disk_autoexpand` - Object that represents autoresize properties of the cluster. It has following attributes:
    * `autoexpand` - Boolean field that indicates whether autoresize is enabled.
    * `max_disk_size` - Maximum disk size for autoresize.

* `wal_volume` - Object that represents wal volume of the cluster. Changing this creates a new cluster. It has following attributes:
    * `size` - (Required) Size of the instance wal volume.
    * `volume_type` - (Required) The type of the cluster wal volume. Changing this creates a new cluster.
    * `autoexpand` - Boolean field that indicates whether wal volume autoresize is enabled.
    * `max_disk_size` - Maximum disk size for wal volume autoresize.

* `network` -  Object that represents network of the cluster. Changing this creates a new cluster. It has following attributes:
    * `uuid` - The id of the network. Changing this creates a new cluster.
    * `port` - The port id of the network. Changing this creates a new cluster.

* `root_enabled` - Boolean field that indicates whether root user is enabled for the cluster.

* `root_password` - Password for the root user of the cluster.

* `configuration_id` - The id of the configuration attached to cluster.

* `capabilities` - Object that represents capability applied to cluster. There can be several instances of this object. Each instance of this object has following attributes:
    * `name` - (Required) The name of the capability to apply.
    * `settings` - Map of key-value settings of the capability.

* `restore_point` - Object that represents backup to restore cluster from. **New since v.0.1.4**. It has following attributes:
    * `backup_id` - (Required) ID of the backup.
    * `target` - Used only for restoring from PITR backups. Timestamp of needed backup in format "2021-10-06 01:02:00". You can specify "latest" to use most recent backup. 

* `backup_schedule` - Object that represents configuration of PITR backup. This functionality is available only for postgres datastore. **New since v.0.1.4** This object has following attributes:
    * `name` - Name of the schedule.
    * `start_hours` - Hours part of timestamp of initial backup
    * `start_minutes` - Minutes part of timestamp of initial backup
    * `interval_hours` - Time interval between backups, specified in hours. Available values: 3, 6, 8, 12, 24.
    * `keep_count` - Number of backups to be stored.


## Import

Clusters can be imported using the `id`, e.g.

```
$ terraform import vkcs_db_cluster.mycluster 708a74a1-6b00-4a96-938c-28a8a6d98590
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.

You should at least add following fields to your .tf file:

`name, flavor_id, cluster_size, volume_size, volume_type, datastore`

Please, use `"IMPORTED"` as value for `volume_type` field.
