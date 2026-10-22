---
subcategory: "Key Management Service (KMS)"
layout: "vkcs"
page_title: "vkcs: vkcs_kms_keys"
description: |-
  List KMS keys within VKCS.
---

# vkcs_kms_keys

Data source containing list of KMS keys

This data source has no configurable arguments. `keys` and `keys_count` are computed when the data source is read.

## Example Usage

```terraform
data "vkcs_kms_keys" "keys" {
  depends_on = [vkcs_kms_key.keys_key_0, vkcs_kms_key.keys_key_1]
}

output "kms_keys" {
  value = data.vkcs_kms_keys.keys.keys
}
```

## Argument Reference

This data source has no configurable arguments.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `keys` *list of* *string* &rarr; Names of the KMS keys. An empty response produces an empty list.

- `keys_count` *number* &rarr; Number of keys in the returned list.

- `id` *string* &rarr; Hexadecimal SHA-256 hash derived from the first ten returned entries (or all entries if fewer than ten) and the total entry count.
