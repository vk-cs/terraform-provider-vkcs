---
subcategory: "Manage Access (IAM)"
layout: "vkcs"
page_title: "vkcs: vkcs_iam_service_user"
description: |-
  Get information on a VKCS IAM service user.
---

# vkcs_iam_service_user

Use this data source to get information about an IAM service user.

## Example Usage

```terraform
data "vkcs_iam_service_user" "service_user" {
  name = vkcs_iam_service_user.service_user.name
}
```

## Argument Reference
- `id` optional *string* &rarr;  ID of the service user.

- `name` optional *string* &rarr;  Name of the service user.

- `region` optional *string* &rarr;  The region in which to obtain the IAM Service Users client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Service user creation timestamp.

- `creator_name` *string* &rarr;  Name of the user who created the service user.

- `description` *string* &rarr;  Service user description.

- `role_names` *string* &rarr;  Names of roles assigned to the service user.


