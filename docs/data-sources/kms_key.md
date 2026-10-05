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
```

## Argument Reference

- `name` **required** *string* &rarr; Name of the KMS key to read.

- `type` optional *string* &rarr; Type of the key, populated from KMS when the data source is read. See the `type` parameter in the [OpenBao transit API documentation](https://openbao.org/docs/api/secret/transit/#parameters). This field does not filter the lookup.

- `deletion_allowed` optional *boolean* &rarr; Whether KMS allows the key to be deleted, populated from KMS when the data source is read. This field does not change the deletion policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` *string* &rarr; Name of the KMS key.
