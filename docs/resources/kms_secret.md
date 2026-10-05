---
subcategory: "Key Management Service (KMS)"
layout: "vkcs"
page_title: "vkcs: vkcs_kms_secret"
description: |-
  Represents a KMS secret resource within VKCS.
---

# vkcs_kms_secret

Resource representing KMS secret

~> **Implementation status:** Creation and updates are not implemented yet, and deletion is not functional with the current schema. The configuration below describes the resource schema; applying it does not create a secret. Use the `vkcs_kms_secret` data source to read an existing secret.

Secret values read into `data` and `data_json` are stored in Terraform state. Neither field is currently marked as sensitive in the resource schema.

## Example Usage

```terraform
resource "vkcs_kms_secret" "secret" {
  path = "application/config"

  data = {
    username = "example-user"
  }
}
```

## Argument Reference

- `path` **required** *string* &rarr; Identifier (path) of the secret within the KMS secret store. Changing this forces a new resource.

- `data` optional computed *map of* *string* &rarr; Data stored in the secret as a Terraform object of string key-value pairs. Refreshed from KMS when the resource is read.

- `data_json` optional *string* &rarr; Data stored in the secret as a JSON object, represented as a string. When the resource is read, this field is populated from the same data as `data`.

- `created_time` optional computed *string* &rarr; Creation time returned in the secret metadata, formatted as a Go time string, for example `2026-09-01 10:00:00 +0000 UTC`.

- `version` optional computed *number* &rarr; Current version of the secret returned in its metadata.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` *string* &rarr; Path of the KMS secret.

## Import

Import is not currently supported.
