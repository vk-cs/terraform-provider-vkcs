---
layout: "vkcs"
page_title: "vkcs: db_cluster_with_shards"
subcategory: ""
description: |-
  Manages a db cluster with shards.
---

# vkcs\_db\_cluster\_with\_shards (Resource)

Provides a db cluster with shards resource. This can be used to create, modify and delete db cluster with shards for clickhouse datastore.

## Example Usage

```terraform

resource "vkcs_db_cluster_with_shards" "db-cluster-with-shards" {
  name = "db-cluster-with-shards"

  datastore {
    type    = "clickhouse"
    version = "20.8"
  }

  shard {
    size        = 2
    shard_id    = example_shard_id1
    flavor_id   = example_flavor_id

    volume_size = 10
    volume_type = example_volume_type
    
    network {
      uuid = example_network_id
    }
  }

  shard {
    size        = 2
    shard_id    = example_shard_id2
    flavor_id   = example_flavor_id
    
    volume_size = 10
    volume_type = example_volume_type

    network {
      uuid = example_network_id
    }
  }
}
```

## Argument Reference

* `name` - (Required) The name of the cluster. Changing this creates a new cluster

* `datastore` - (Required) Object that represents datastore of the cluster. Changing this creates a new cluster. It has following attributes:
    * `type` - (Required) Type of the datastore. Changing this creates a new cluster. Type of the datastore must be "clickhouse".
    * `version` - (Required) Version of the datastore. Changing this creates a new cluster.

* `keypair` - Name of the keypair to be attached to cluster. Changing this creates a new cluster.

* `floating_ip_enabled` - Boolean field that indicates whether floating ip is created for cluster. Changing this creates a new cluster.

* `root_enabled` - Boolean field that indicates whether root user is enabled for the cluster.

* `root_password` - Password for the root user of the cluster.

* `configuration_id` - The id of the configuration attached to cluster.

* `capabilities` - Object that represents capability applied to cluster. There can be several instances of this object. Each instance of this object has following attributes:
    * `name` - (Required) The name of the capability to apply.
    * `settings` - Map of key-value settings of the capability.

* `shard` - (Required) Object that represents cluster shard. There can be several instances of this object. Each instance of this object has following attributes:
    * `size` - (Required) The number of instances in the cluster shard.
    * `shard_id` - (Required) The ID of the shard. Changing this creates a new cluster.
    * `flavor_id` - (Required) The ID of flavor for the cluster shard.
    * `availability_zone` - The name of the availability zone of the cluster shard. Changing this creates a new cluster.
    * `volume_size` - (Required) Size of the cluster shard instance volume.
    * `volume_type` - (Required) The type of the cluster shard instance volume.
    * `wal_volume` - Object that represents wal volume of the cluster. It has following attributes:
        * `size` - (Required) Size of the instance wal volume.
        * `volume_type` - (Required) The type of the cluster wal volume.
        * `autoexpand` - Boolean field that indicates whether wal volume autoresize is enabled.
        * `max_disk_size` - Maximum disk size for wal volume autoresize.
    * `network` -  Object that represents network of the cluster shard. Changing this creates a new cluster. It has following attributes: 
        * `uuid` - The id of the network. Changing this creates a new cluster.
        * `port` - The port id of the network. Changing this creates a new cluster.

## Import

Clusters can be imported using the `id`, e.g.

```
$ terraform import vkcs_db_cluster_with_shards.mycluster 708a74a1-6b00-4a96-938c-28a8a6d98590
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.

You should at least add following fields to your .tf file:

`name, datastore`, and for each shard add: `shard_id, size, flavor_id, volume_size, volume_type`

Please, use `"IMPORTED"` as value for `volume_type` field.