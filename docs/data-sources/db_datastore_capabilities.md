---
subcategory: "Databases"
layout: "vkcs"
page_title: "vkcs: vkcs_db_datastore_capabilities"
description: |-
  Get information on capabilities supported for a VKCS db datastore.
---

# vkcs_db_datastore_capabilities

Use this data source to get capabilities supported for a VKCS datastore.

**New since v0.2.0**.

## Example Usage

```terraform
data "vkcs_db_datastore_capabilities" "postgres_caps" {
  datastore_name       = data.vkcs_db_datastore.postgres
  datastore_version_id = local.pg_v14_version_id
}

output "postgresql_capabilities" {
  value       = data.vkcs_db_datastore_capabilities.postgres_caps.capabilities
  description = "Available capabilities of the latest version of PostgreSQL datastore."
}
```

## Argument Reference
- `datastore_name` **required** *string* &rarr;  Name of the data store.

- `datastore_version_id` **required** *string* &rarr;  ID of the version of the data store.

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `capabilities`  *list* &rarr;  Capabilities of the datastore.
  - `allow_major_upgrade` *boolean* &rarr;  This attribute indicates whether a capability can be applied in the next major version of data store.

  - `allow_upgrade_from_backup` *boolean* &rarr;  This attribute indicates whether a capability can be applied to upgrade from backup.

  - `description` *string* &rarr;  Description of data store capability.

  - `name` *string* &rarr;  Name of data store capability.

  - `params`  *list*
    - `default_value` *string* &rarr;  Default value for a parameter.

    - `element_type` *string* &rarr;  Type of element value for a parameter of `list` type.

    - `enum_values` *string* &rarr;  Supported values for a parameter.

    - `masked` *boolean* &rarr;  Masked indicates whether a parameter value must be a boolean mask.

    - `max` *number* &rarr;  Maximum value for a parameter.

    - `min` *number* &rarr;  Minimum value for a parameter.

    - `name` *string* &rarr;  Name of a parameter.

    - `regex` *string* &rarr;  Regular expression that a parameter value must match.

    - `required` *boolean* &rarr;  Required indicates whether a parameter value must be set.

    - `type` *string* &rarr;  Type of value for a parameter.


  - `should_be_on_master` *boolean* &rarr;  This attribute indicates whether a capability applies only to the master node.


- `id` *string* &rarr;  ID of the resource


