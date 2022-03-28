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

```terraform

resource "vkcs_db_cluster" "mydb-cluster" {
  name        = "mydb-cluster"

  datastore {
    type    = "postgresql"
    version = "12"
  }

  cluster_size = 3

  flavor_id   = example_flavor_id

  volume_size = 10
  volume_type = example_volume_type

  network {
    uuid = example_network_id
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

## Import

Clusters can be imported using the `id`, e.g.

```
$ terraform import vkcs_db_cluster.mycluster 708a74a1-6b00-4a96-938c-28a8a6d98590
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.

You should at least add following fields to your .tf file:

`name, flavor_id, cluster_size, volume_size, volume_type, datastore`

Please, use `"IMPORTED"` as value for `volume_type` field.
