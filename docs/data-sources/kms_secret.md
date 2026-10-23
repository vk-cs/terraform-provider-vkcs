---
subcategory: "Key Management Service (KMS)"
layout: "vkcs"
page_title: "vkcs: vkcs_kms_secret"
description: |-
  Get information about a KMS secret within VKCS.
---

# vkcs_kms_secret

A data source containing a KMS secret

Secret values in `data_json` are marked as sensitive and are stored in Terraform state.

Only `path` is configurable. `data_json`, `created_time`, `version`, and `delete_protection` are computed from KMS when the data source is read.

## Example Usage

```terraform
data "vkcs_kms_secret" "secret" {
  path = vkcs_kms_secret.secret.id
}

output "password" {
  value     = jsondecode(data.vkcs_kms_secret.secret.data_json).password
  sensitive = true
}
```

## Argument Reference

- `path` **required** *string* &rarr; Identifier (path) of the secret within the KMS secret store.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `data_json` sensitive *string* &rarr; Secret data as a JSON object, represented as a string and populated from KMS when the data source is read. Use `jsondecode` to access its fields.

- `created_time` *string* &rarr; Creation time returned in the secret metadata, formatted as a Go time string, for example `2026-09-01 10:00:00 +0000 UTC`.

- `version` *number* &rarr; Current version of the secret returned in its metadata.

- `delete_protection` *boolean* &rarr; Whether the secret is protected from deletion.

- `id` *string* &rarr; Path of the KMS secret.
