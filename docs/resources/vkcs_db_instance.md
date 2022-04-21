---
layout: "vkcs"
page_title: "vkcs: db_instance"
subcategory: ""
description: |-
  Manages a db instance.
---

# vkcs\_db\_instance

Provides a db instance resource. This can be used to create, modify and delete db instance.

## Example Usage

```terraform

resource "vkcs_db_instance" "db-instance" {
  name = "db-instance"

  datastore {
    type    = example_datastore_type
    version = example_datastore_version
  }

  floating_ip_enabled = true

  flavor_id         = example_flavor_id
  availability_zone = example_availability_zone

  size        = 8
  volume_type = example_volume_type
  disk_autoexpand {
    autoexpand    = true
    max_disk_size = 1000
  }

  network {
    uuid = example_network_id
  }

  capabilities {
    name = capability_name
  }

  capabilities {
    name = another_capability_name
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the instance. Changing this creates a new instance

* `replica_of` - ID of the instance, that current instance is replica of.

* `datastore` - (Required) Object that represents datastore of the instance. Changing this creates a new instance. It has following attributes:
    * `type` - (Required) Type of the datastore. Changing this creates a new instance.
    * `version` - (Required) Version of the datastore. Changing this creates a new instance.

* `keypair` - Name of the keypair to be attached to instance. Changing this creates a new instance.

* `floating_ip_enabled` - Boolean field that indicates whether floating ip is created for instance. Changing this creates a new instance.

* `flavor_id` - (Required) The ID of flavor for the instance.

* `availability_zone` - The name of the availability zone of the instance. Changing this creates a new instance.

* `size` - (Required) Size of the instance volume.

* `volume_type` - (Required) The type of the instance volume. Changing this creates a new instance.

* `disk_autoexpand` - Object that represents autoresize properties of the instance. It has following attributes:
    * `autoexpand` - Boolean field that indicates whether autoresize is enabled.
    * `max_disk_size` - Maximum disk size for autoresize.
  
* `wal_disk_autoexpand` - Object that represents autoresize properties of the instance. It has following attributes:
    * `autoexpand` - Boolean field that indicates whether wal volume autoresize is enabled.
    * `max_disk_size` - Maximum disk size for wal volume autoresize.

* `wal_volume` - Object that represents wal volume of the instance. Changing this creates a new instance. It has following attributes:
    * `size` - (Required) Size of the instance wal volume.
    * `volume_type` - (Required) The type of the instance wal volume. Changing this creates a new instance.

* `network` -  Object that represents network of the instance. Changing this creates a new instance. It has following attributes: 
    * `uuid` - The id of the network. Changing this creates a new instance.
    * `port` - The port id of the network. Changing this creates a new instance.
    * `fixed_ip_v4` - The IPv4 address. Changing this creates a new instance.

* `root_enabled` - Boolean field that indicates whether root user is enabled for the instance.

* `root_password` - Password for the root user of the instance. If this field is empty and root user is enabled, then after creation of the instance this field will contain auto-generated root user password.

* `configuration_id` - The id of the configuration attached to instance.

* `capabilities` - Object that represents capability applied to instance. There can be several instances of this object (see example). Each instance of this object has following attributes:
    * `name` - (Required) The name of the capability to apply.
    * `settings` - Map of key-value settings of the capability.

## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import vkcs_db_instance.myinstance 708a74a1-6b00-4a96-938c-28a8a6d98590
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.

You should at least add following fields to your .tf file:

`name, flavor_id, size, volume_type, datastore`

Please, use `"IMPORTED"` as value for `volume_type` field.