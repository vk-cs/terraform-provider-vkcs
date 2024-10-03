---
subcategory: "CDN"
layout: "vkcs"
page_title: "vkcs: vkcs_cdn_shielding_pop"
description: |-
  Get information on a VKCS CDN shielding point of presence (POP).
---

# vkcs_cdn_shielding_pop



## Example Usage

```terraform
data "vkcs_cdn_shielding_pop" "pop" {
    city = "Moscow-Megafon"
}
```

## Argument Reference
- `city` optional *string* &rarr;  City of origin shielding location.

- `country` optional *string* &rarr;  Country of origin shielding location.

- `datacenter` optional *string* &rarr;  Name of origin shielding location datacenter.

- `region` optional *string* &rarr;  The region in which to obtain the CDN client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *number* &rarr;  ID of the origin shielding location.


