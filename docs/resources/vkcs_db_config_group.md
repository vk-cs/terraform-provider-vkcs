---
layout: "vkcs"
page_title: "vkcs: vkcs_db_config_group"
description: |-
  Manages a db config group.
---

# vkcs_db_config_group

Provides a db config group resource. This can be used to create, update and delete db config group.
**New since v.0.1.7**.

## Example Usage

```terraform
resource "vkcs_db_config_group" "db-config-group" {
    name = "db-config-group"
    datastore {
        type = "mysql"
        version = "8.0"
    }
    values = {
        activate_all_roles_on_login : "true"
        autocommit : "1"
        block_encryption_mode : "test"
        innodb_segment_reserve_factor : "0.53"
    }
    description = "db-config-group-description"
}


resource "vkcs_db_instance" "db-instance" {
    name = "db-instance"

    availability_zone = "GZ1"
    
    datastore {
        type = "mysql"
        version = "8.0"
    }
    
    configuration_id = vkcs_db_config_group.db-config-group.id
    network {
      uuid = vkcs_networking_network.db.id
    }
    flavor_id = data.vkcs_compute_flavor.db.id
    volume_type = "ceph-ssd"
    size = 8

    depends_on = [
        vkcs_networking_router_interface.db
    ]
}
```
## Argument Reference
- `datastore` (***Required***) Object that represents datastore of the config group. Changing this creates a new config group.
  - `type` **String** (***Required***) Type of the datastore.

  - `version` **String** (***Required***) Version of the datastore.

- `name` **String** (***Required***) The name of the config group.

- `values` <strong>Map of </strong>**String** (***Required***) Map of configuration parameters in format "key": "value".

- `description` **String** (*Optional*) The description of the config group.


## Attributes Reference
- `datastore`  See Argument Reference above.
  - `type` **String** See Argument Reference above.

  - `version` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `values` <strong>Map of </strong>**String** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `created` **String** Timestamp of config group's creation

- `id` **String** ID of the resource.

- `updated` **String** Timestamp of config group's last update



## Updating config group

While it is possible to create/delete config groups that are not attached to any instance or cluster, in order to update config group, it must be attached.

## Import

Config groups can be imported using the `id`, e.g.

```shell
terraform import vkcs_db_config_group.myconfiggroup d3d6f037-84f6-44f7-a9f4-ac4b40d67859
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.
