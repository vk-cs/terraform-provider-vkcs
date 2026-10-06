---
subcategory: "Key Management Service (KMS)"
layout: "vkcs"
page_title: "vkcs: vkcs_kms_key"
description: |-
  Represents a KMS key resource within VKCS.
---

# vkcs_kms_key

Resource representing KMS key

Keys are protected from deletion by default. Before destroying a key, set `deletion_allowed = true` and apply the change.

## Example Usage

```terraform
resource "vkcs_kms_key" "key" {
  name             = "key-tf-example"
  type             = "aes256-gcm96"
  deletion_allowed = false
}
```

## Argument Reference

- `name` **required** *string* &rarr; Name of the KMS key. Changing this forces a new resource.

- `type` optional *string* &rarr; Type of the KMS key: `aes128-gcm96`, `aes256-gcm96`, `chacha20-poly1305`, or `xchacha20-poly1305`. Defaults to `aes256-gcm96`. Changing this forces a new resource.

- `deletion_allowed` optional *boolean* &rarr; Whether KMS allows the key to be deleted. Defaults to `false`. Set to `true` and apply before destroying the key. Refreshed from KMS when the resource is read.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` *string* &rarr; Name of the KMS key.

## Import

Import is not currently supported.
