---
layout: "vkcs"
page_title: "Managing CDN resources"
description: |-
  Managing CDN resources with VKCS Provider
---

# Manage CDN with the VKCS Terraform Provider

This guide provides a detailed, step-by-step approach to configuring and managing VKCS Content Delivery Network (CDN) resources using the provider. It covers essential tasks such as setting up CDN resources, origin groups, and SSL certificates.

## Prerequisites

Before diving into the guide, ensure you meet the following prerequisites:

- **Configure Terraform and VKCS Provider** Make sure that you installed Terraform CLI and configured VKCS Provider. Follow [instructions](https://registry.terraform.io/providers/vk-cs/vkcs/latest/docs/guides/getting_started) if needed.
- **Understand Terraform Basics:** Familiarize yourself with Terraform concepts like resource lifecycles, dependencies, and state management. [Terraform documentation](https://developer.hashicorp.com/terraform/docs) will help you understand the basic principles and key points.
- **Understand VKCS CDN Basics:** Check the CDN service [documentation](https://cloud.vk.com/docs/en/networks/cdn) to understand main concepts.

## Configuring Origin Groups

Origin groups manage the backend servers responsible for hosting and delivering content. A well-configured origin group ensures reliability and efficiency, and you should always configure an one for a CDN resource with an use of `vkcs_cdn_origin_group` resource.

Consider adding multiple origins for redundancy and failover. To get more details, refer to the VK Cloud [documentation](https://cloud.vk.com/docs/en/networks/cdn/concepts/origin-groups).

### Example Configuration:

```terraform
resource "vkcs_cdn_origin_group" "origin_group" {
  name = "tfexample-origin-group"
  origins = [
    {
      source = "origin1.vk.com"
    },
    {
      source = "origin2.vk.com",
      backup = true
    }
  ]
  use_next = true
}
```

## Adding a SSL certificate

Proceed to the next step if you plan to use Let's Encrypt certificate, or not use one at all.

SSL certificates should be used for content delivery over HTTPS protocol. To manage your own certificates and to apply them to CDN resources, you can use `vkcs_cdn_ssl_certificate` resource.

```terraform
resource "vkcs_cdn_ssl_certificate" "certificate" {
  name        = "tfexample-ssl-certificate"
  certificate = file("${path.module}/certificate.pem")
  private_key = file("${path.module}/private-key.key")
}
```

!> **Security Note:** `certificate` and `private_key` are marked as sensitive, and, therefore, will not be shown in outputs, but you should consider protecting them as input variables and state values. To get more information on the topic, you can refer to the [official tutorial](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables).

##  Utilizing Shielding PoPs

Proceed to the next step if you do not plan to enable CDN resource shielding or if it is not available in your region.

Shielding PoPs act as intermediaries to enhance performance by caching content closer to end users, you should choose one strategically based on traffic origins.

### List all Points of Presence

To list all points of presence, you can use "vkcs_cdn_shielding_pops" data source.

```terraform
data "vkcs_cdn_shielding_pops" "pops" {}

output "shielding_locations" {
  value = data.vkcs_cdn_shielding_pops.pops.shielding_pops
}
```

### Retrieve the identifier of PoP

To enable shielding on a CDN resource, you should provide the identifier of a specific PoP, which can be retrieved with an use of `vkcs_cdn_shielding_pop` data source:

```terraform
data "vkcs_cdn_shielding_pop" "pop" {
  city = "Moscow-Megafon"
}
```

## Creating a CDN Resource

CDN resources serve as the cornerstone for accelerating content delivery, optimizing reliability, and reducing the load of origin services. To create an one, you should use `vkcs_cdn_resource` resource.

### Example Configuration:

```terraform
resource "vkcs_cdn_resource" "resource" {
  cname        = local.cname # Provide your own value
  origin_group = vkcs_cdn_origin_group.origin_group.id
  options = {
    edge_cache_settings = {
      value = "10m"
    }
    forward_host_header = true
  }
  # Remove if you decided not to enable shielding on the resource
  shielding = {
    enabled = true
    pop_id  = data.vkcs_cdn_shielding_pop.pop.id
  }
  # Remove if not necessary. Check provider's documentation for
  # the attribute to get more information on how to provide a SSL
  # certificate for a CDN resource.
  ssl_certificate = {
    type = "own"
    id   = vkcs_cdn_ssl_certificate.certificate.id
  }
}
```

## Complete Example Configuration

Below is an integrated example showcasing all components:

```terraform
resource "vkcs_cdn_origin_group" "origin_group" {
  name = "tfexample-origin-group"
  origins = [
    {
      source = "origin1.vk.com"
    },
    {
      source = "origin2.vk.com",
      backup = true
    }
  ]
  use_next = true
}

# Remove if not needed
resource "vkcs_cdn_ssl_certificate" "certificate" {
  name        = "tfexample-ssl-certificate"
  certificate = file("${path.module}/certificate.pem")
  private_key = file("${path.module}/private-key.key")
}

# Remove if not needed
data "vkcs_cdn_shielding_pop" "pop" {
  city = "Moscow-Megafon"
}

resource "vkcs_cdn_resource" "resource" {
  cname        = local.cname # Provide your own value
  origin_group = vkcs_cdn_origin_group.origin_group.id
  options = {
    edge_cache_settings = {
      value = "10m"
    }
    forward_host_header = true
  }
  # Remove if you decided not to enable shielding on the resource
  shielding = {
    enabled = true
    pop_id  = data.vkcs_cdn_shielding_pop.pop.id
  }
  # Remove if not necessary. Check provider's documentation for
  # the attribute to get more information on how to provide a SSL
  # certificate for a CDN resource.
  ssl_certificate = {
    type = "own"
    id   = vkcs_cdn_ssl_certificate.certificate.id
  }
}
```

## Next Steps

Review the full documentation on CDN management with the VKCS Terraform Provider in the corresponding category of the provider [documentation](https://registry.terraform.io/providers/vk-cs/vkcs/latest/docs), pay special attention on the available CDN resource options. Test various configurations to optimize content freshness and perfomance, and to customize access.
