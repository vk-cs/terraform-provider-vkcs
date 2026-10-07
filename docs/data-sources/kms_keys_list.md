---
subcategory: "Key Management Service (KMS)"
layout: "vkcs"
page_title: "vkcs: vkcs_kms_keys_list"
description: |-
  List KMS keys within VKCS.
---

# vkcs_kms_keys_list

Data source containing list of KMS keys

## Example Usage

```terraform
resource "vkcs_kms_key" "keys_list_key_0" {
  name             = "keys_list_key_0"
  type             = "aes128-gcm96"
  deletion_allowed = true
}

resource "vkcs_kms_key" "keys_list_key_1" {
  name             = "keys_list_key_1"
  deletion_allowed = true
}

data "vkcs_kms_keys_list" "keys" {
  depends_on = [vkcs_kms_key.keys_list_key_0, vkcs_kms_key.keys_list_key_1]
}

output "kms_keys" {
  value = data.vkcs_kms_keys_list.keys.keys
}
```

## Argument Reference

No arguments are required. The following optional computed fields are populated when the data source is read; they do not filter the results.

- `keys` optional computed *list of* *string* &rarr; Names of the KMS keys. An empty response produces an empty list.

- `keys_count` optional computed *number* &rarr; Number of keys in the returned list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` *string* &rarr; Hexadecimal SHA-256 hash derived from the first ten returned entries (or all entries if fewer than ten) and the total entry count.
