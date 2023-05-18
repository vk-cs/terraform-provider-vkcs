---
subcategory: "Databases"
layout: "vkcs"
page_title: "vkcs: vkcs_db_datastores"
description: |-
  Get information on db datastores that are available within VKCS.
---

# vkcs_db_datastores

Use this data source to get a list of datastores from VKCS. **New since v0.2.0**.

## Example Usage

```terraform
data "vkcs_db_datastores" "datastores" {}

output "available_datastores" {
    value = data.vkcs_db_datastores.datastores.datastores
    description = "List of datastores that are available within VKCS."
}
```

## Argument Reference
- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `datastores` *object*

- `id` *string* &rarr;  ID of the resource.


