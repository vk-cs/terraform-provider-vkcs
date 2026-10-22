---
subcategory: "Key Management Service (KMS)"
layout: "vkcs"
page_title: "vkcs: vkcs_kms_secrets"
description: |-
  List KMS secrets within VKCS.
---

# vkcs_kms_secrets

Data source containing list of KMS secrets

This data source has no configurable arguments. `secrets` and `secrets_count` are computed when the data source is read.

## Example Usage

```terraform
data "vkcs_kms_secrets" "secrets" {
  depends_on = [vkcs_kms_secret.secrets_secret_0, vkcs_kms_secret.secrets_secret_1]
}

output "kms_secrets" {
  value = data.vkcs_kms_secrets.secrets.secrets
}
```

## Argument Reference

This data source has no configurable arguments.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `secrets` *list of* *string* &rarr; Paths of the secrets returned by KMS. An empty response produces an empty list.

- `secrets_count` *number* &rarr; Number of secrets in the returned list.

- `id` *string* &rarr; Hexadecimal SHA-256 hash derived from the first ten returned entries (or all entries if fewer than ten) and the total entry count.
