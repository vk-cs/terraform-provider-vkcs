---
layout: "vkcs"
page_title: "vkcs: db_config_group"
subcategory: ""
description: |-
  Get information on a db config group.
---

# vkcs\_db\_config_group

Use this data source to get the information on a db config group resource.
**New since v.0.1.7**.

## Example Usage

```terraform

data "vkcs_db_config_group" "db-config-group" {
  config_group_id = "7a914e84-8fcf-46f8-bbe5-a8337ba090f4"
}
```
## Argument Reference

The following arguments are supported:

* `config_group_id` - (Required) The UUID of the config_group.

## Attributes reference

The following attributes are exported:

* `name` - The name of the config group.

* `datastore` - Object that represents datastore of backup
    * `type` - Type of the datastore.
    * `version` - Version of the datastore.

* `description` - The description of the config group.

* `values` - Map of configuration parameters in format "key": "value".  

* `created` - Timestamp of config group's creation

* `updated` - Timestamp of config group's last update