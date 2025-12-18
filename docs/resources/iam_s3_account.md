---
subcategory: "Manage Access (IAM)"
layout: "vkcs"
page_title: "vkcs: vkcs_iam_s3_account"
description: |-
  Manages an IAM S3 account resource within VKCS.
---

# vkcs_iam_s3_account

Manages an IAM S3 account within VKCS.

!> **Security Note:** `secret_key` is marked as sensitive, and, therefore, will not be shown in outputs by default, but you should consider protecting it as input variable and state value. To get more information on the topic, you can refer to the [official tutorial](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables).

## Example Usage
```terraform
resource "vkcs_iam_s3_account" "s3_account" {
  name        = "tf-example-s3-account"
  description = "S3 account created by Terraform example"
}

output "access_key" {
  value = vkcs_iam_s3_account.s3_account.access_key
}

output "secret_key" {
  value     = vkcs_iam_s3_account.s3_account.secret_key
  sensitive = true
}
```

## Argument Reference
- `name` **required** *string* &rarr;  Name of the S3 account. The name must be unique. The length must be between 3 and 32 characters. Changing this creates a new resource.

- `description` optional *string* &rarr;  Description of the S3 account. Changing this creates a new resource.

- `region` optional *string* &rarr;  The region in which to obtain the IAM Service Users client. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `access_key` *string* &rarr;  Access key for the S3 account.

- `account_id` *string* &rarr;  ID of the S3 account in Hotbox S3 service.

- `account_name` *string* &rarr;  Name of the S3 account in Hotbox S3 service.

- `created_at` *string* &rarr;  S3 account creation timestamp.

- `id` *string* &rarr;  ID of the S3 account.

- `secret_key` *string* &rarr;  Secret key for the S3 account. <br>**Note:** This is a sensitive attribute.



## Import

An IAM S3 account can be imported using the `id`, e.g.
```shell
terraform import vkcs_iam_s3_account.s3_account <s3_account_id>
```
