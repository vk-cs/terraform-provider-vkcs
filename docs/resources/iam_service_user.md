---
subcategory: "Manage Access (IAM)"
layout: "vkcs"
page_title: "vkcs: vkcs_iam_service_user"
description: |-
  Manages an IAM service user resource within VKCS.
---

# vkcs_iam_service_user

Manages an IAM service user within VKCS.

!> **Security Note:** `password` is marked as sensitive, and, therefore, will not be shown in outputs by default, but you should consider protecting it as input variable and state value. To get more information on the topic, you can refer to the [official tutorial](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables).

To get information on available roles for service accounts, refer to to the VK Cloud [documentation](https://cloud.vk.com/docs/en/tools-for-using-services/account/concepts/rolesandpermissions).

## Example Usage
```terraform
resource "vkcs_iam_service_user" "service_user" {
  name        = "tf-example-service-user"
  description = "Service user created by Terraform example"
  role_names = [
    "mcs_admin_vm",
    "mcs_admin_network"
  ]
}

output "credentials" {
  value = {
    login    = vkcs_iam_service_user.service_user.login
    password = vkcs_iam_service_user.service_user.password
  }
  sensitive = true
}
```

## Argument Reference
- `name` **required** *string* &rarr;  Name of the service user. The name must be unique. The length must be between 3 and 32 characters. Changing this creates a new resource.

- `role_names` **required** *string* &rarr;  Names of roles assigned to the service user. Changing this creates a new resource.

- `description` optional *string* &rarr;  Description of the service user. The maximum length is 256 characters. Changing this creates a new resource.

- `region` optional *string* &rarr;  The region in which to obtain the IAM Service Users client. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Service user creation timestamp.

- `creator_name` *string* &rarr;  Name of the user who created the service user.

- `id` *string* &rarr;  ID of the service user.

- `login` *string* &rarr;  Login name of the service user.

- `password` *string* &rarr;  Password of the service user. <br>**Note:** This is a sensitive attribute.



## Import

An IAM service account can be imported using the `id`, e.g.
```shell
terraform import vkcs_iam_service_user.service_user <service_user_id>
```
