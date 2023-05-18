---
subcategory: "Databases"
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
  cloud_monitoring_enabled = true

  volume_size = 10
  volume_type = "ceph-ssd"

  network {
    uuid = vkcs_networking_network.db.id
  }

  depends_on = [
    vkcs_networking_router_interface.db
  ]
}

data "vkcs_lb_loadbalancer" "loadbalancer" {
  id = "${vkcs_db_cluster.db-cluster.loadbalancer_id}"
}

data "vkcs_networking_port" "loadbalancer-port" {
  port_id = "${data.vkcs_lb_loadbalancer.loadbalancer.vip_port_id}"
}

output "cluster_ips" {
  value = "${data.vkcs_networking_port.loadbalancer-port.all_fixed_ips}"
  description = "IP addresses of the cluster."
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
- `cluster_size` **required** *number* &rarr;  The number of instances in the cluster.

- `datastore` **required** &rarr;  Object that represents datastore of the cluster. Changing this creates a new cluster.
  - `type` **required** *string* &rarr;  Type of the datastore. Changing this creates a new cluster. Type of the datastore can either be "galera_mysql", "postgresql" or "tarantool".

  - `version` **required** *string* &rarr;  Version of the datastore. Changing this creates a new cluster.

- `flavor_id` **required** *string* &rarr;  The ID of flavor for the cluster.

- `name` **required** *string* &rarr;  The name of the cluster. Changing this creates a new cluster.

- `volume_size` **required** *number* &rarr;  Size of the cluster instance volume.

- `volume_type` **required** *string* &rarr;  The type of the cluster instance volume. Changing this creates a new cluster.

- `availability_zone` optional *string* &rarr;  The name of the availability zone of the cluster. Changing this creates a new cluster.

- `backup_schedule` optional &rarr;  Object that represents configuration of PITR backup. This functionality is available only for postgres datastore. **New since v0.1.4**.
  - `interval_hours` **required** *number* &rarr;  Time interval between backups, specified in hours. Available values: 3, 6, 8, 12, 24.

  - `keep_count` **required** *number* &rarr;  Number of backups to be stored.

  - `name` **required** *string* &rarr;  Name of the schedule.

  - `start_hours` **required** *number* &rarr;  Hours part of timestamp of initial backup.

  - `start_minutes` **required** *number* &rarr;  Minutes part of timestamp of initial backup.

- `capabilities` optional &rarr;  Object that represents capability applied to cluster. There can be several instances of this object.
  - `name` **required** *string* &rarr;  The name of the capability to apply.

  - `settings` optional *map of* *string* &rarr;  Map of key-value settings of the capability.

- `cloud_monitoring_enabled` optional *boolean* &rarr;  Enable cloud monitoring for the cluster. Changing this for Redis or MongoDB creates a new instance. **New since v0.2.0**.

- `configuration_id` optional *string* &rarr;  The id of the configuration attached to cluster.

- `disk_autoexpand` optional &rarr;  Object that represents autoresize properties of the cluster.
  - `autoexpand` optional *boolean* &rarr;  Indicates whether autoresize is enabled.

  - `max_disk_size` optional *number* &rarr;  Maximum disk size for autoresize.

- `floating_ip_enabled` optional *boolean* &rarr;  Indicates whether floating ip is created for cluster. Changing this creates a new cluster.

- `keypair` optional *string* &rarr;  Name of the keypair to be attached to cluster. Changing this creates a new cluster.

- `network` optional &rarr;  Object that represents network of the cluster. Changing this creates a new cluster.
  - `port` optional deprecated *string* &rarr;  The port id of the network. Changing this creates a new cluster. **Deprecated** This argument is deprecated, please do not use it.

  - `security_groups` optional *set of* *string* &rarr;  An array of one or more security group IDs to associate with the cluster instances. Changing this creates a new cluster. **New since v0.2.0**.

  - `subnet_id` optional *string* &rarr;  The id of the subnet. Changing this creates a new cluster. **New since v0.1.15**.

  - `uuid` optional *string* &rarr;  The id of the network. Changing this creates a new cluster.**Note** Although this argument is marked as optional, it is actually required at the moment. Not setting a value for it may cause an error.

- `region` optional *string* &rarr;  Region to create resource in.

- `restore_point` optional &rarr;  Object that represents backup to restore cluster from. **New since v0.1.4**.
  - `backup_id` **required** *string* &rarr;  ID of the backup.

  - `target` optional *string* &rarr;  Used only for restoring from PITR backups. Timestamp of needed backup in format "2021-10-06 01:02:00". You can specify "latest" to use most recent backup.

- `root_enabled` optional *boolean* &rarr;  Indicates whether root user is enabled for the cluster.

- `root_password` optional sensitive *string* &rarr;  Password for the root user of the cluster.

- `shrink_options` optional *string* &rarr;  Used only for shrinking cluster. List of IDs of instances that should remain after shrink. If no options are supplied, shrink operation will choose first non-leader instance to delete.

- `wal_disk_autoexpand` optional &rarr;  Object that represents autoresize properties of wal volume of the cluster.
  - `autoexpand` optional *boolean* &rarr;  Indicates whether wal volume autoresize is enabled.

  - `max_disk_size` optional *number* &rarr;  Maximum disk size for wal volume autoresize.

- `wal_volume` optional &rarr;  Object that represents wal volume of the cluster. Changing this creates a new cluster.
  - `size` **required** *number* &rarr;  Size of the instance wal volume.

  - `volume_type` **required** *string* &rarr;  The type of the cluster wal volume. Changing this creates a new cluster.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `instances` *object* &rarr;  Cluster instances info.

- `loadbalancer_id` *string* &rarr;  The id of the loadbalancer attached to the cluster. **New since v0.1.15**.



## Import

Clusters can be imported using the `id`, e.g.

```shell
terraform import vkcs_db_cluster.mycluster 708a74a1-6b00-4a96-938c-28a8a6d98590
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.

You should at least add following fields to your .tf file:

`name, flavor_id, cluster_size, volume_size, volume_type, datastore`

Please, use `"IMPORTED"` as value for `volume_type` field.
