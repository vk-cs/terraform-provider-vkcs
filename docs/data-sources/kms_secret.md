---
subcategory: "Key Management Service (KMS)"
layout: "vkcs"
page_title: "vkcs: vkcs_kms_secret"
description: |-
  Get information about a KMS secret within VKCS.
---

# vkcs_kms_secret

A data source containing a KMS secret

Secret values read into `data` and `data_json` are stored in Terraform state. Neither field is currently marked as sensitive in the data source schema.

## Example Usage

```terraform
data "vkcs_kms_secret" "secret" {
  path = "application/config"
}

output "secret_data" {
  value     = data.vkcs_kms_secret.secret.data
  sensitive = true
}
```

## Argument Reference

- `path` **required** *string* &rarr; Identifier (path) of the secret within the KMS secret store.

- `data` optional computed *map of* *string* &rarr; Data stored in the secret as a Terraform object of string key-value pairs, populated from KMS when the data source is read.

- `data_json` optional *string* &rarr; Data stored in the secret as a JSON object, represented as a string containing the same key-value pairs returned in `data`.

- `created_time` optional computed *string* &rarr; Creation time returned in the secret metadata, formatted as a Go time string, for example `2026-09-01 10:00:00 +0000 UTC`.

- `version` optional computed *number* &rarr; Current version of the secret returned in its metadata. This field does not select a version to read.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` *string* &rarr; Path of the KMS secret.
