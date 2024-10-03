---
subcategory: "CDN"
layout: "vkcs"
page_title: "vkcs: vkcs_cdn_shielding_pops"
description: |-
  Get information on available VKCS CDN shielding points of presence (POPs).
---

# vkcs_cdn_shielding_pops



## Example Usage

```terraform
data "vkcs_cdn_shielding_pops" "pops" {}

output "shielding_locations" {
  value = data.vkcs_cdn_shielding_pops.pops.shielding_pops
}
```

## Argument Reference
- `region` optional *string* &rarr;  The region in which to obtain the CDN client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `shielding_pops`  *list* &rarr;  List of available origin shielding points of precense.
  - `city` *string* &rarr;  City of origin shielding location.

  - `country` *string* &rarr;  Country of origin shielding location.

  - `datacenter` *string* &rarr;  Name of origin shielding location datacenter.

  - `id` *number* &rarr;  ID of the origin shielding location.



