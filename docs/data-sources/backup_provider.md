---
subcategory: "Backup"
layout: "vkcs"
page_title: "vkcs: vkcs_backup_provider"
description: |-
  Get information on an VKCS backup provider.
---

# vkcs_backup_provider

Use this data source to get backup provider info

**New since v0.4.0**.

## Example Usage

```terraform
data "vkcs_backup_provider" "cloud-servers" {
  name = "cloud_servers"
}
```

## Argument Reference
- `name` **required** *string* &rarr;  Name of the backup provider

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource


