---
subcategory: "CDN"
layout: "vkcs"
page_title: "vkcs: vkcs_cdn_ssl_certificate"
description: |-
  Manages a CDN SSL certificate within VKCS.
---

# vkcs_cdn_ssl_certificate



## Example Usage
```terraform
resource "vkcs_cdn_ssl_certificate" "certificate" {
  name        = "tfexample-ssl-certificate"
  certificate = file("${path.module}/certificate.pem")
  private_key = file("${path.module}/private-key.key")
}
```

## Argument Reference
- `certificate` **required** sensitive *string* &rarr;  Public part of the SSL certificate. All chain of the SSL certificate should be added.

- `name` **required** *string* &rarr;  SSL certificate name.

- `private_key` **required** sensitive *string* &rarr;  Private key of the SSL certificate.

- `region` optional *string* &rarr;  The region in which to obtain the CDN client. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *number* &rarr;  ID of the SSL certificate.

- `issuer` *string* &rarr;  Name of the certification center issued the SSL certificate.

- `subject_cn` *string* &rarr;  Domain name that the SSL certificate secures.

- `validity_not_after` *string* &rarr;  Date when certificate become untrusted (ISO 8601/RFC 3339 format, UTC.).

- `validity_not_before` *string* &rarr;  Date when certificate become valid (ISO 8601/RFC 3339 format, UTC.).



## Import

A SSL certificate can be imported using the `id`, e.g.
```shell
terraform import vkcs_cdn_resource.resource <resource_id>
```
