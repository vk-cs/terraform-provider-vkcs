---
subcategory: "Key Management Service (KMS)"
layout: "vkcs"
page_title: "vkcs: vkcs_kms_secret"
description: |-
  Represents a KMS secret resource within VKCS.
---

# vkcs_kms_secret

Resource representing KMS secret

Secrets are protected from deletion by default. Before destroying a secret, set `delete_protection = false` and apply the change.

Secret values in `data_json` are marked as sensitive and are stored in Terraform state.

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
```

## Argument Reference

- `path` **required** *string* &rarr; Identifier (path) of the secret within the KMS secret store. Changing this forces a new resource.

- `data_json` **required** *string* &rarr; Data stored in the secret as a JSON object, represented as a string. Use `jsonencode` to encode a Terraform object.

- `created_time` optional computed *string* &rarr; Creation time returned in the secret metadata, formatted as a Go time string, for example `2026-09-01 10:00:00 +0000 UTC`.

- `version` optional computed *number* &rarr; Current version of the secret returned in its metadata.

- `delete_protection` optional *boolean* &rarr; Whether the secret is protected from deletion. Defaults to `true`. Set to `false` and apply before destroying the secret.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` *string* &rarr; Path of the KMS secret.

## Import

Import is not currently supported.
