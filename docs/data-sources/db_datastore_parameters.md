---
subcategory: "Databases"
layout: "vkcs"
page_title: "vkcs: vkcs_db_datastore_parameters"
description: |-
  Get information on configuration parameters supported for a VKCS db datastore.
---

# vkcs_db_datastore_parameters

Use this data source to get configuration parameters supported for a VKCS datastore. **New since v0.2.0**.

## Example Usage

```terraform
data "vkcs_db_datastore_parameters" "mysql_params" {
	datastore_name = data.vkcs_db_datastore.mysql
	datastore_version_id = local.mysql_v8_version_id
}

output "mysql_parameters" {
	value = data.vkcs_db_datastore_parameters.mysql_params.parameters
	description = "Available configuration parameters of the latest version of MySQL datastore."
}
```

## Argument Reference
- `datastore_name` **required** *string* &rarr;  Name of the data store.

- `datastore_version_id` **required** *string* &rarr;  ID of the version of the data store.

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `parameters` *object* &rarr;  Versions of the datastore.


