---
subcategory: "Key Management Service (KMS)"
layout: "vkcs"
page_title: "vkcs: vkcs_kms_key"
description: |-
  Get information about a KMS key within VKCS.
---

# vkcs_kms_key

A data source containing a KMS key

## Example Usage

```terraform
data "vkcs_kms_key" "key" {
  name = "key-tf-example"
}

output "key_type" {
  value = data.vkcs_kms_key.key.type
}

output "key_deletion_allowed" {
  value = data.vkcs_kms_key.key.deletion_allowed
}
```

## Argument Reference

- `name` **required** *string* &rarr; Name of the KMS key to read.

- `type` optional *string* &rarr; Type of the key: `aes128-gcm96`, `aes256-gcm96`, `chacha20-poly1305`, or `xchacha20-poly1305`. Populated from KMS when the data source is read. This field does not filter the lookup.

- `deletion_allowed` optional *boolean* &rarr; Whether KMS allows the key to be deleted, populated from KMS when the data source is read. This field does not change the deletion policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` *string* &rarr; Name of the KMS key.
