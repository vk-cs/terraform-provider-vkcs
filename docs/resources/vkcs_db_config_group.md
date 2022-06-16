---
layout: "vkcs"
page_title: "vkcs: db_config_group"
subcategory: ""
description: |-
  Manages a db config group.
---

# vkcs\_db\_config_group

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
    description = "test-desc"
}


resource "vkcs_db_instance" "db-instance" {
    name = "db-instance"
    datastore {
        type = "mysql"
        version = "8.0"
    }

    configuration_id = vkcs_db_config_group.db-config-group.id
    network {
      uuid = "403321eb-b939-49dc-b2b4-a802beb74222"
    }
    flavor_id = "c8c42890-1ae9-411f-8cce-42e2d7c9b7d0"
    volume_type = "ceph-ssd"
    size = 8
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the config group.

* `datastore` - (Required) Object that represents datastore of the config group. Changing this creates a new config group. It has following attributes:
    * `type` - (Required) Type of the datastore. Changing this creates a new config group.
    * `version` - (Required) Version of the datastore. Changing this creates a new config group.

* `description` - The description of the config group.

* `values` - Map of configuration parameters in format "key": "value".  

## Attributes reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `datastore` - See Argument Reference above.

* `description` - See Argument Reference above.

* `values` - See Argument Reference above.

* `created` - Timestamp of config group's creation

* `updated` - Timestamp of config group's last update

## Updating config group

While it is possible to create/delete config groups that are not attached to any instance or cluster, in order to update config group, it must be attached.

## Import

Config groups can be imported using the `id`, e.g.

```
$ terraform import vkcs_db_config_group.myconfiggroup d3d6f037-84f6-44f7-a9f4-ac4b40d67859
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.
