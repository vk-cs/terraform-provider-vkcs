---
subcategory: "Manage Access (IAM)"
layout: "vkcs"
page_title: "vkcs: vkcs_iam_s3_account"
description: |-
  Get information on a VKCS IAM S3 account.
---

# vkcs_iam_s3_account

Use this data source to get information about an IAM S3 account.

## Example Usage

```terraform
data "vkcs_iam_s3_account" "s3_account" {
  name = vkcs_iam_s3_account.s3_account.name
}
```

## Argument Reference
- `id` optional *string* &rarr;  ID of the S3 account. Conflicts with `name`.

- `name` optional *string* &rarr;  Name of the S3 account. Conflicts with `id`.

- `region` optional *string* &rarr;  The region in which to obtain the IAM Service Users client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `access_key` *string* &rarr;  Access key for the S3 account.

- `account_id` *string* &rarr;  ID of the S3 account in Hotbox S3 service.

- `account_name` *string* &rarr;  Name of the S3 account in Hotbox S3 service.

- `created_at` *string* &rarr;  S3 account creation timestamp.

- `description` *string* &rarr;  Description of the S3 account.


