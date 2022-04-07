---
layout: "vkcs"
page_title: "Getting Started with the VKCS Provider"
description: |-
  Getting Started with the VKCS Provider
---

# Create basic config for VKCS Provider resources

This example shows how to create a simple terraform configuration for creation of VKCS resources.

First, create a Terraform config file named `main.tf`. Inside, you'll want to include the configuration of
[VKCS Provider](https://registry.terraform.io/providers/vk-cs/vkcs/latest/docs).

Use VKCS provider:

```hcl
provider "vkcs" {
    username   = "some_user"
    password   = "s3cr3t"
    project_id = "some_project_id"
  }
}
```

Configure VKCS provider:

* The `username` field should be replaced with your user_name.
* The `password` field should be replaced with your user's password.
* The `project_id` field should be replaced with your project_id.

For additional configuration parameters, please read [configuration reference](https://registry.terraform.io/providers/vk-cs/vkcs/latest/docs#configuration-reference)
