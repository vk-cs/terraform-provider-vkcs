---
subcategory: "Key Management Service (KMS)"
layout: "vkcs"
page_title: "vkcs: vkcs_kms_key"
description: |-
  Represents a KMS key resource within VKCS.
---

# vkcs_kms_key

Resource representing KMS key

~> **Implementation status:** Creation, updates, and deletion are not implemented yet. The configuration below describes the resource schema; applying it does not create a key. Use the `vkcs_kms_key` data source to read an existing key.

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

- `type` optional *string* &rarr; Type of the KMS key. See the `type` parameter in the [OpenBao transit API documentation](https://openbao.org/docs/api/secret/transit/#parameters). Refreshed from KMS when the resource is read. Changing this forces a new resource.

- `deletion_allowed` optional *boolean* &rarr; Whether KMS allows the key to be deleted. Refreshed from KMS when the resource is read.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` *string* &rarr; Name of the KMS key.

## Import

Import is not currently supported.
