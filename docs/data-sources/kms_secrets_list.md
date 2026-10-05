---
subcategory: "Key Management Service (KMS)"
layout: "vkcs"
page_title: "vkcs: vkcs_kms_secrets_list"
description: |-
  List KMS secrets within VKCS.
---

# vkcs_kms_secrets_list

Data source containing list of KMS secrets

## Example Usage

```terraform
data "vkcs_kms_secrets_list" "secrets" {}

output "kms_secrets" {
  value = data.vkcs_kms_secrets_list.secrets.secrets
}
```

## Argument Reference

No arguments are required. The following optional computed fields are populated when the data source is read; they do not filter the results.

- `secrets` optional computed *list of* *string* &rarr; Paths of the secrets returned by KMS. An empty response produces an empty list.

- `secrets_count` optional computed *number* &rarr; Number of secrets in the returned list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` *string* &rarr; Hexadecimal SHA-256 hash derived from the first ten returned entries (or all entries if fewer than ten) and the total entry count.
