---
layout: "vkcs"
page_title: "Provider: VKCS"
description: |-
  The VKCS provider is used to interact with VKCS services.
  The provider needs to be configured with the proper credentials before it can be used.
---

# VKCS Provider

Test changes

The VKCS provider is used to interact with
[VKCS services](https://mcs.mail.ru/). The provider needs
to be configured with the proper credentials before it can be used.

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

## Configuration Reference

The following arguments are supported:

* `username` - (Required) The username to login with.
  If omitted, the `USER_NAME` environment variable is used.

* `password` - (Required) The Password to login with. If omitted, the `PASSWORD` environment variable is used.

* `project_id` - (Required) The ID of Project to login with. 
  If omitted, the `PROJECT_ID` environment variable is used.

* `auth_url` - (Optional) URL for authentication in VKCS. Default is https://infra.mail.ru/identity/v3/.

* `region` - (Optional) A region to use. Default is `RegionOne`.

