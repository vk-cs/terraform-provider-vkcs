---
layout: "vkcs"
page_title: "vkcs: vkcs_db_config_group"
description: |-
  Get information on a db config group.
---

# vkcs_db_config_group

Use this data source to get the information on a db config group resource.
**New since v.0.1.7**.

## Example Usage

```terraform
data "vkcs_db_config_group" "db-config-group" {
  config_group_id = "7a914e84-8fcf-46f8-bbe5-a8337ba090f4"
}
```

## Argument Reference
- `config_group_id` **String** (***Required***) The UUID of the config_group.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created` **String** Timestamp of config group's creation.

- `datastore` **Object** Object that represents datastore of backup

- `description` **String** The description of the config group.

- `id` **String** ID of the resource.

- `name` **String** The name of the config group.

- `updated` **String** Timestamp of config group's last update.

- `values` <strong>Map of </strong>**String** Map of configuration parameters in format "key": "value".


