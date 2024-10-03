---
subcategory: "CDN"
layout: "vkcs"
page_title: "vkcs: vkcs_cdn_ssl_certificate"
description: |-
  Get information on a VKSC CDN SSL certificate.
---

# vkcs_cdn_ssl_certificate



## Example Usage

```terraform
data "vkcs_cdn_ssl_certificate" "cert" {
  name = vkcs_cdn_ssl_certificate.certificate.name
}
```

## Argument Reference
- `name` **required** *string* &rarr;  SSL certificate name.

- `region` optional *string* &rarr;  The region in which to obtain the CDN client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *number* &rarr;  ID of the SSL certificate.

- `issuer` *string* &rarr;  Name of the certification center issued the SSL certificate.

- `subject_cn` *string* &rarr;  Domain name that the SSL certificate secures.

- `validity_not_after` *string* &rarr;  Date when certificate become untrusted (ISO 8601/RFC 3339 format, UTC.).

- `validity_not_before` *string* &rarr;  Date when certificate become valid (ISO 8601/RFC 3339 format, UTC.).


