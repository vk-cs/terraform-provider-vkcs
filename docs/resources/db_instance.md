---
subcategory: "Databases"
layout: "vkcs"
page_title: "vkcs: vkcs_db_instance"
description: |-
  Manages a db instance.
---

# vkcs_db_instance

Provides a db instance resource. This can be used to create, modify and delete db instance.

## Example Usage
### Basic instance
```terraform
resource "vkcs_db_instance" "db-instance" {
  name        = "db-instance"

  availability_zone = "GZ1"

  datastore {
    type    = "mysql"
    version = "5.7"
  }

  flavor_id   = data.vkcs_compute_flavor.db.id
  cloud_monitoring_enabled = true

  size        = 8
  volume_type = "ceph-ssd"
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }

  network {
    uuid = vkcs_networking_network.db.id
    security_groups = [vkcs_networking_secgroup.secgroup.id]
  }

  capabilities {
    name = "node_exporter"
    settings = {
      "listen_port" : "9100"
    }
  }

  depends_on = [
    vkcs_networking_router_interface.db,
    vkcs_networking_secgroup.secgroup
  ]
}
```

### Instance restored from backup
```terraform
resource "vkcs_db_instance" "db-instance" {
  name = "db-instance"

  datastore {
    type    = "postgresql"
    version = "13"
  }

  floating_ip_enabled = true

  flavor_id         = "9e931469-1490-489e-88af-29a289681c53"
  availability_zone = "MS1"

  size        = 8
  volume_type = "MS1"
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }

  network {
    uuid = "3ee9b184-3311-4d85-840b-7a9c48e7beac"
  }

  capabilities {
    name = "node_exporter"
  }

  capabilities {
    name = "postgres_extensions"
  }

  restore_point {
    backup_id = "backup_id"
  }
}
```

### Postgresql instance with scheduled PITR backup
```terraform
resource "vkcs_db_instance" "db-instance" {
  name = "db-instance"

  datastore {
    type    = "postgresql"
    version = "13"
  }

  floating_ip_enabled = true

  flavor_id         = "9e931469-1490-489e-88af-29a289681c53"
  availability_zone = "MS1"

  size        = 8
  volume_type = "MS1"
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }

  network {
    uuid = "3ee9b184-3311-4d85-840b-7a9c48e7beac"
  }

  capabilities {
    name = "node_exporter"
  }

  capabilities {
    name = "postgres_extensions"
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
- `datastore` **required** &rarr;  Object that represents datastore of the instance. Changing this creates a new instance.
  - `type` **required** *string* &rarr;  Type of the datastore. Changing this creates a new instance.

  - `version` **required** *string* &rarr;  Version of the datastore. Changing this creates a new instance.

- `flavor_id` **required** *string* &rarr;  The ID of flavor for the instance.

- `name` **required** *string* &rarr;  The name of the instance. Changing this creates a new instance

- `size` **required** *number* &rarr;  Size of the instance volume.

- `volume_type` **required** *string* &rarr;  The type of the instance volume. Changing this creates a new instance.

- `availability_zone` optional *string* &rarr;  The name of the availability zone of the instance. Changing this creates a new instance.

- `backup_schedule` optional &rarr;  Object that represents configuration of PITR backup. This functionality is available only for postgres datastore. **New since v0.1.4**.
  - `interval_hours` **required** *number* &rarr;  Time interval between backups, specified in hours. Available values: 3, 6, 8, 12, 24.

  - `keep_count` **required** *number* &rarr;  Number of backups to be stored.

  - `name` **required** *string* &rarr;  Name of the schedule.

  - `start_hours` **required** *number* &rarr;  Hours part of timestamp of initial backup.

  - `start_minutes` **required** *number* &rarr;  Minutes part of timestamp of initial backup.

- `capabilities` optional &rarr;  Object that represents capability applied to instance. There can be several instances of this object (see example).
  - `name` **required** *string* &rarr;  The name of the capability to apply.

  - `settings` optional *map of* *string* &rarr;  Map of key-value settings of the capability.

- `cloud_monitoring_enabled` optional *boolean* &rarr;  Enable cloud monitoring for the instance. Changing this for Redis or MongoDB creates a new instance. **New since v0.2.0**.

- `configuration_id` optional *string* &rarr;  The id of the configuration attached to instance.

- `disk_autoexpand` optional &rarr;  Object that represents autoresize properties of the instance.
  - `autoexpand` optional *boolean* &rarr;  Indicates whether autoresize is enabled.

  - `max_disk_size` optional *number* &rarr;  Maximum disk size for autoresize.

- `floating_ip_enabled` optional *boolean* &rarr;  Indicates whether floating ip is created for instance. Changing this creates a new instance.

- `keypair` optional *string* &rarr;  Name of the keypair to be attached to instance. Changing this creates a new instance.

- `network` optional &rarr;  Object that represents network of the instance. Changing this creates a new instance.
  - `fixed_ip_v4` optional *string* &rarr;  The IPv4 address. Changing this creates a new instance. **Note** This argument conflicts with "replica_of". Setting both at the same time causes "fixed_ip_v4" to be ignored.

  - `port` optional deprecated *string* &rarr;  The port id of the network. Changing this creates a new instance. **Deprecated** This argument is deprecated, please do not use it.

  - `security_groups` optional *set of* *string* &rarr;  An array of one or more security group IDs to associate with the instance. Changing this creates a new instance. **New since v0.2.0**.

  - `subnet_id` optional *string* &rarr;  The id of the subnet. Changing this creates a new instance. **New since v0.1.15**.

  - `uuid` optional *string* &rarr;  The id of the network. Changing this creates a new instance.**Note** Although this argument is marked as optional, it is actually required at the moment. Not setting a value for it may cause an error.

- `region` optional *string* &rarr;  Region to create resource in.

- `replica_of` optional *string* &rarr;  ID of the instance, that current instance is replica of.

- `restore_point` optional &rarr;  Object that represents backup to restore instance from. **New since v0.1.4**.
  - `backup_id` **required** *string* &rarr;  ID of the backup.

  - `target` optional *string* &rarr;  Used only for restoring from postgresql PITR backups. Timestamp of needed backup in format "2021-10-06 01:02:00". You can specify "latest" to use most recent backup.

- `root_enabled` optional *boolean* &rarr;  Indicates whether root user is enabled for the instance.

- `root_password` optional sensitive *string* &rarr;  Password for the root user of the instance. If this field is empty and root user is enabled, then after creation of the instance this field will contain auto-generated root user password.

- `wal_disk_autoexpand` optional &rarr;  Object that represents autoresize properties of the instance wal volume.
  - `autoexpand` optional *boolean* &rarr;  Indicates whether wal volume autoresize is enabled.

  - `max_disk_size` optional *number* &rarr;  Maximum disk size for wal volume autoresize.

- `wal_volume` optional &rarr;  Object that represents wal volume of the instance. Changing this creates a new instance.
  - `size` **required** *number* &rarr;  Size of the instance wal volume.

  - `volume_type` **required** *string* &rarr;  The type of the instance wal volume.

  - `autoexpand` optional deprecated *boolean* &rarr;  Indicates whether wal volume autoresize is enabled. **Deprecated** Please, use wal_disk_autoexpand block instead.

  - `max_disk_size` optional deprecated *number* &rarr;  Maximum disk size for wal volume autoresize. **Deprecated** Please, use wal_disk_autoexpand block instead.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `ip` *string* &rarr;  IP address of the instance.



## Import

Instances can be imported using the `id`, e.g.

```shell
terraform import vkcs_db_instance.myinstance 708a74a1-6b00-4a96-938c-28a8a6d98590
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.

You should at least add following fields to your .tf file:

`name, flavor_id, size, volume_type, datastore`

Please, use `"IMPORTED"` as value for `volume_type` field.
