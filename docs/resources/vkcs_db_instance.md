---
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
  
  size        = 8
  volume_type = "ceph-ssd"
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }

  network {
    uuid = vkcs_networking_network.db.id
  }

  capabilities {
    name = "node_exporter"
    settings = {
      "listen_port" : "9100"
    }
  }

  depends_on = [
    vkcs_networking_router_interface.db
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
- `datastore` (***Required***) Object that represents datastore of the instance. Changing this creates a new instance.
  - `type` **String** (***Required***) Type of the datastore. Changing this creates a new instance.

  - `version` **String** (***Required***) Version of the datastore. Changing this creates a new instance.

- `flavor_id` **String** (***Required***) The ID of flavor for the instance.

- `name` **String** (***Required***) The name of the instance. Changing this creates a new instance

- `size` **Number** (***Required***) Size of the instance volume.

- `volume_type` **String** (***Required***) The type of the instance volume. Changing this creates a new instance.

- `availability_zone` **String** (*Optional*) The name of the availability zone of the instance. Changing this creates a new instance.

- `backup_schedule` (*Optional*) Object that represents configuration of PITR backup. This functionality is available only for postgres datastore. **New since v.0.1.4**.
  - `interval_hours` **Number** (***Required***) Time interval between backups, specified in hours. Available values: 3, 6, 8, 12, 24.

  - `keep_count` **Number** (***Required***) Number of backups to be stored.

  - `name` **String** (***Required***) Name of the schedule.

  - `start_hours` **Number** (***Required***) Hours part of timestamp of initial backup.

  - `start_minutes` **Number** (***Required***) Minutes part of timestamp of initial backup.

- `capabilities` (*Optional*) Object that represents capability applied to instance. There can be several instances of this object (see example).
  - `name` **String** (***Required***) The name of the capability to apply.

  - `settings` <strong>Map of </strong>**String** (*Optional*) Map of key-value settings of the capability.

- `configuration_id` **String** (*Optional*) The id of the configuration attached to instance.

- `disk_autoexpand` (*Optional*) Object that represents autoresize properties of the instance.
  - `autoexpand` **Boolean** (*Optional*) Indicates whether autoresize is enabled.

  - `max_disk_size` **Number** (*Optional*) Maximum disk size for autoresize.

- `floating_ip_enabled` **Boolean** (*Optional*) Indicates whether floating ip is created for instance. Changing this creates a new instance.

- `ip` **String** (*Optional*) IP address of the instance.

- `keypair` **String** (*Optional*) Name of the keypair to be attached to instance. Changing this creates a new instance.

- `network` (*Optional*) Object that represents network of the instance. Changing this creates a new instance.
  - `fixed_ip_v4` **String** (*Optional*) The IPv4 address. Changing this creates a new instance. **Note** This argument conflicts with "replica_of". Setting both at the same time causes "fixed_ip_v4" to be ignored.

  - `port` **String** (*Optional* Deprecated) The port id of the network. Changing this creates a new instance. ***Deprecated*** This argument is deprecated, please do not use it.

  - `subnet_id` **String** (*Optional*) The id of the subnet. Changing this creates a new instance. **New since v.0.1.15**.

  - `uuid` **String** (*Optional*) The id of the network. Changing this creates a new instance.**Note** Although this argument is marked as optional, it is actually required at the moment. Not setting a value for it may cause an error.

- `region` **String** (*Optional*) Region to create resource in.

- `replica_of` **String** (*Optional*) ID of the instance, that current instance is replica of.

- `restore_point` (*Optional*) Object that represents backup to restore instance from. **New since v.0.1.4**.
  - `backup_id` **String** (***Required***) ID of the backup.

  - `target` **String** (*Optional*) Used only for restoring from postgresql PITR backups. Timestamp of needed backup in format "2021-10-06 01:02:00". You can specify "latest" to use most recent backup.

- `root_enabled` **Boolean** (*Optional*) Indicates whether root user is enabled for the instance.

- `root_password` **String** (*Optional* Sensitive) Password for the root user of the instance. If this field is empty and root user is enabled, then after creation of the instance this field will contain auto-generated root user password.

- `wal_disk_autoexpand` (*Optional*) Object that represents autoresize properties of the instance wal volume.
  - `autoexpand` **Boolean** (*Optional*) Indicates whether wal volume autoresize is enabled.

  - `max_disk_size` **Number** (*Optional*) Maximum disk size for wal volume autoresize.

- `wal_volume` (*Optional*) Object that represents wal volume of the instance. Changing this creates a new instance.
  - `size` **Number** (***Required***) Size of the instance wal volume.

  - `volume_type` **String** (***Required***) The type of the instance wal volume.

  - `autoexpand` **Boolean** (*Optional* Deprecated) Indicates whether wal volume autoresize is enabled. ***Deprecated***. Please, use wal_disk_autoexpand block instead.

  - `max_disk_size` **Number** (*Optional* Deprecated) Maximum disk size for wal volume autoresize. ***Deprecated***. Please, use wal_disk_autoexpand block instead.


## Attributes Reference
- `datastore`  See Argument Reference above.
  - `type` **String** See Argument Reference above.

  - `version` **String** See Argument Reference above.

- `flavor_id` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `size` **Number** See Argument Reference above.

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

- `ip` **String** See Argument Reference above.

- `keypair` **String** See Argument Reference above.

- `network`  See Argument Reference above.
  - `fixed_ip_v4` **String** See Argument Reference above.

  - `port` **String** See Argument Reference above.

  - `subnet_id` **String** See Argument Reference above.

  - `uuid` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `replica_of` **String** See Argument Reference above.

- `restore_point`  See Argument Reference above.
  - `backup_id` **String** See Argument Reference above.

  - `target` **String** See Argument Reference above.

- `root_enabled` **Boolean** See Argument Reference above.

- `root_password` **String** See Argument Reference above.

- `wal_disk_autoexpand`  See Argument Reference above.
  - `autoexpand` **Boolean** See Argument Reference above.

  - `max_disk_size` **Number** See Argument Reference above.

- `wal_volume`  See Argument Reference above.
  - `size` **Number** See Argument Reference above.

  - `volume_type` **String** See Argument Reference above.

  - `autoexpand` **Boolean** See Argument Reference above.

  - `max_disk_size` **Number** See Argument Reference above.

- `id` **String** ID of the resource.



## Import

Instances can be imported using the `id`, e.g.

```shell
terraform import vkcs_db_instance.myinstance 708a74a1-6b00-4a96-938c-28a8a6d98590
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.

You should at least add following fields to your .tf file:

`name, flavor_id, size, volume_type, datastore`

Please, use `"IMPORTED"` as value for `volume_type` field.
