---
layout: "vkcs"
page_title: "Provider: VKCS Provider"
description: |-
  The VKCS provider is used to interact with VKCS services.
  The provider needs to be configured with the proper credentials before it can be used.
---

# VKCS Provider

The VKCS provider is used to interact with [VKCS services](https://mcs.mail.ru/). The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

Terraform 1.0 and later:

```terraform
# Configure the vkcs provider

terraform {
  required_providers {
    vkcs = {
      source = "vk-cs/vkcs"
      version = "~> 0.1.0"
    }
  }
}

# Create new compute instance
resource "vkcs_compute_instance" "myinstance"{
  # ...
}
```

## Authentication

The VKCS provider supports username/password authentication. Preconfigured provider file with `username` and `project_id` can be downloaded from [https://mcs.mail.ru/app/project](https://mcs.mail.ru/app/project) portal. Go to `Terraform` tab -> click on the "Download VKCS provider file".

```terraform
provider "vkcs" {
    username   = "USERNAME"
    password   = "PASSWORD"
    project_id = "PROJECT_ID"
}
```

## Argument Reference
- `auth_url` optional *string* &rarr;  The Identity authentication URL.

- `cloud_containers_api_version` optional *string* &rarr;  Cloud Containers API version to use.
_NOTE_ Only for custom VKCS deployments.

- `password` optional sensitive *string* &rarr;  Password to login with.

- `project_id` optional *string* &rarr;  The ID of Project to login with.

- `region` optional *string* &rarr;  A region to use.

- `user_domain_id` optional *string* &rarr;  The id of the domain where the user resides.

- `user_domain_name` optional *string* &rarr;  The name of the domain where the user resides.

- `username` optional *string* &rarr;  User name to login with.



