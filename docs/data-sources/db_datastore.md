---
subcategory: "Databases"
layout: "vkcs"
page_title: "vkcs: vkcs_db_datastore"
description: |-
  Get information on a VKCS db datastore.
---

# vkcs_db_datastore

Use this data source to get information on a VKCS db datastore. **New since v0.2.0**.

## Example Usage

```terraform
data "vkcs_db_datastore" "datastore" {
  name = "mysql"
}

output "mysql_versions" {
  value = data.vkcs_db_datastore.datastore.versions
  description = "List of versions of MySQL that are available within VKCS."
}
```

## Argument Reference
- `id` optional *string* &rarr;  The id of the datastore.

- `name` optional *string* &rarr;  The name of the datastore.

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `cluster_volume_types` *string* &rarr;  Supported volume types for the datastore when used in a cluster.

- `minimum_cpu` *number* &rarr;  Minimum CPU required for instance of the datastore.

- `minimum_ram` *number* &rarr;  Minimum RAM required for instance of the datastore.

- `versions` *object* &rarr;  Versions of the datastore.

- `volume_types` *string* &rarr;  Supported volume types for the datastore.


