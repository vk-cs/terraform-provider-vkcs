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

Terraform 1.1.5 and later:

```terraform
# Configure the vkcs provider

terraform {
  required_providers {
    vkcs = {
      source = "vk-cs/vkcs"
      version = "< 1.0.0"
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
- `access_token` optional sensitive *string* &rarr;  A temporary token to use for authentication. You alternatively can use `OS_AUTH_TOKEN` environment variable. If both are specified, this attribute takes precedence. <br>**Note:** The token will not be renewed and will eventually expire, usually after 1 hour. If access is needed for longer than a token's lifetime, use credentials-based authentication.

- `auth_url` optional *string* &rarr;  The Identity authentication URL.

- `cloud_containers_api_version` optional *string* &rarr;  Cloud Containers API version to use. <br>**Note:** Only for custom VKCS deployments.

- `endpoint_overrides` optional &rarr;  Custom endpoints for corresponding APIs. If not specified, endpoints provided by the catalog will be used.
    - `backup` optional *string* &rarr;  Backup API custom endpoint.

    - `block_storage` optional *string* &rarr;  Block Storage API custom endpoint.

    - `cdn` optional *string* &rarr;  CDN API custom endpoint.

    - `compute` optional *string* &rarr;  Compute API custom endpoint.

    - `container_infra` optional *string* &rarr;  Cloud Containers API custom endpoint.

    - `container_infra_addons` optional *string* &rarr;  Cloud Containers Addons API custom endpoint.

    - `data_platform` optional *string* &rarr;  Data Platform API custom endpoint.

    - `database` optional *string* &rarr;  Database API custom endpoint.

    - `iam_service_users` optional *string* &rarr;  IAM Service Users API custom endpoint.

    - `ics` optional *string* &rarr;  ICS API custom endpoint.

    - `image` optional *string* &rarr;  Image API custom endpoint.

    - `key_manager` optional *string* &rarr;  Key Manager API custom endpoint.

    - `load_balancer` optional *string* &rarr;  Load Balancer API custom endpoint.

    - `ml_platform` optional *string* &rarr;  ML Platform API custom endpoint.

    - `networking` optional *string* &rarr;  Networking API custom endpoint.

    - `public_dns` optional *string* &rarr;  Public DNS API custom endpoint.

    - `shared_filesystem` optional *string* &rarr;  Shared Filesystem API custom endpoint.

    - `templater` optional *string* &rarr;  Templater API custom endpoint.

- `password` optional sensitive *string* &rarr;  Password to login with.

- `project_id` optional *string* &rarr;  The ID of Project to login with.

- `region` optional *string* &rarr;  A region to use.

- `skip_client_auth` optional *boolean* &rarr;  Skip authentication on client initialization. Only applicablie if `access_token` is provided. <br>**Note:** If set to true, the endpoint catalog will not be used for discovery and all required endpoints must be provided via `endpoint_overrides`.

- `user_domain_id` optional *string* &rarr;  The id of the domain where the user resides.

- `user_domain_name` optional *string* &rarr;  The name of the domain where the user resides.

- `username` optional *string* &rarr;  User name to login with.




## Working with VKCS Cloud Storage

VKCS provider does not support working with cloud storage.
To do this, we recommend an AWS provider, you can learn how to use it by following this [documentation](https://cloud.vk.com/docs/en/tools-for-using-services/terraform/how-to-guides/aws).
