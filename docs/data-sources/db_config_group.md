---
subcategory: "Databases"
layout: "vkcs"
page_title: "vkcs: vkcs_db_config_group"
description: |-
  Get information on a db config group.
---

# vkcs_db_config_group

Use this data source to get the information on a db config group resource.

**New since v0.1.7**.

## Example Usage

```terraform
data "vkcs_db_config_group" "db-config-group" {
  config_group_id = "7a914e84-8fcf-46f8-bbe5-a8337ba090f4"
}
```

## Argument Reference
- `config_group_id` **required** *string* &rarr;  The UUID of the config_group.

- `region` optional *string* &rarr;  The region in which to obtain the service client. If omitted, the `region` argument of the provider is used.<br>**New since v0.4.0**.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created` *string* &rarr;  Timestamp of config group's creation.

- `datastore`  *list* &rarr;  Object that represents datastore of backup
  - `type` *string* &rarr;  Type of the datastore.

  - `version` *string* &rarr;  Version of the datastore.


- `description` *string* &rarr;  The description of the config group.

- `id` *string* &rarr;  ID of the resource.

- `name` *string* &rarr;  The name of the config group.

- `updated` *string* &rarr;  Timestamp of config group's last update.

- `values` *map of* *string* &rarr;  Map of configuration parameters in format "key": "value".


