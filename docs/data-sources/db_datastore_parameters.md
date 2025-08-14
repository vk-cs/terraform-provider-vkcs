---
subcategory: "Databases"
layout: "vkcs"
page_title: "vkcs: vkcs_db_datastore_parameters"
description: |-
  Get information on configuration parameters supported for a VKCS db datastore.
---

# vkcs_db_datastore_parameters

Use this data source to get configuration parameters supported for a VKCS datastore.

**New since v0.2.0**.

## Example Usage

```terraform
data "vkcs_db_datastore_parameters" "mysql_params" {
  datastore_name       = data.vkcs_db_datastore.mysql
  datastore_version_id = local.mysql_v8_version_id
}

output "mysql_parameters" {
  value       = data.vkcs_db_datastore_parameters.mysql_params.parameters
  description = "Available configuration parameters of the latest version of MySQL datastore."
}
```

## Argument Reference
- `datastore_name` **required** *string* &rarr;  Name of the data store.

- `datastore_version_id` **required** *string* &rarr;  ID of the version of the data store.

- `region` optional *string* &rarr;  The region to obtain the service client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `parameters`  *list* &rarr;  Configuration parameters supported for the datastore.
    - `max` *number* &rarr;  Maximum value of a configuration parameter.

    - `min` *number* &rarr;  Minimum value of a configuration parameter.

    - `name` *string* &rarr;  Name of a configuration parameter.

    - `restart_required` *boolean* &rarr;  This attribute indicates whether a restart required when a parameter is set.

    - `type` *string* &rarr;  Type of a configuration parameter.



