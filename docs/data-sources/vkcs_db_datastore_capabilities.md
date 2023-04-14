---
layout: "vkcs"
page_title: "vkcs: vkcs_db_datastore_capabilities"
description: |-
  Get information on capabilities supported for a VKCS db datastore.
---

# vkcs_db_datastore_capabilities

Use this data source to get capabilities supported for a VKCS datastore. **New since v.0.2.0**.

## Example Usage

```terraform
data "vkcs_db_datastore_capabilities" "postgres_caps" {
	datastore_name = data.vkcs_db_datastore.postgres
	datastore_version_id = local.pg_v14_version_id
}

output "postgresql_capabilities" {
	value = data.vkcs_db_datastore_capabilities.postgres_caps.capabilities
	description = "Available capabilities of the latest version of PostgreSQL datastore."
}
```

## Argument Reference
- `datastore_name` **String** (***Required***) Name of the data store.

- `datastore_version_id` **String** (***Required***) ID of the version of the data store.

- `region` **String** (*Optional*) The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
- `datastore_name` **String** See Argument Reference above.

- `datastore_version_id` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `capabilities` **Object** Versions of the datastore.

- `id` **String** ID of the resource.

