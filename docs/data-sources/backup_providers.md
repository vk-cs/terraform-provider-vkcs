---
subcategory: "Backup"
layout: "vkcs"
page_title: "vkcs: vkcs_backup_providers"
description: |-
  Get information on an VKCS backup providers.
---

# vkcs_backup_providers

Use this data source to get backup providers info

**New since v0.4.0**.

## Example Usage

```terraform
data "vkcs_backup_providers" "providers" {}
```

## Argument Reference
- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `providers`  *list* &rarr;  List of available backup providers
    - `name` *string* &rarr;  Name of the backup provider

    - `id` *string* &rarr;  ID of the backup provider



