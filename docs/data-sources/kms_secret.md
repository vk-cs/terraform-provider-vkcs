---
subcategory: "Key Management Service (KMS)"
layout: "vkcs"
page_title: "vkcs: vkcs_kms_secret"
description: |-
  Get information about a KMS secret within VKCS.
---

# vkcs_kms_secret

A data source containing a KMS secret

Secret values in `data_json` are stored in Terraform state. This field is not marked as sensitive in the schema.

## Example Usage

```terraform
resource "vkcs_kms_secret" "secret" {
  path              = "credentials"
  delete_protection = false

  data_json = jsonencode({
    username = "my-username"
    password = "my-password"
  })
}

data "vkcs_kms_secret" "secret" {
  path = vkcs_kms_secret.secret.id
}

output "secret_data" {
  value     = jsondecode(data.vkcs_kms_secret.secret.data_json)
  sensitive = true
}
```

## Argument Reference

- `path` **required** *string* &rarr; Identifier (path) of the secret within the KMS secret store.

- `data_json` optional *string* &rarr; Secret data as a JSON object, represented as a string and populated from KMS when the data source is read. Use `jsondecode` to access its fields.

- `created_time` optional computed *string* &rarr; Creation time returned in the secret metadata, formatted as a Go time string, for example `2026-09-01 10:00:00 +0000 UTC`.

- `version` optional computed *number* &rarr; Current version of the secret returned in its metadata. This field does not select a version to read.

- `delete_protection` optional *boolean* &rarr; Whether the secret is protected from deletion, populated from KMS when the data source is read. This field does not change the deletion policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` *string* &rarr; Path of the KMS secret.
