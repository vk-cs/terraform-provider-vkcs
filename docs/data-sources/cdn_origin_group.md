---
subcategory: "CDN"
layout: "vkcs"
page_title: "vkcs: vkcs_cdn_origin_group"
description: |-
  Get information on a VKCS CDN origin group.
---

# vkcs_cdn_origin_group



## Example Usage

```terraform
data "vkcs_cdn_origin_group" "origin_group" {
  name = vkcs_cdn_origin_group.origin_group.name
}
```

## Argument Reference
- `name` **required** *string* &rarr;  Name of the origin group.

- `region` optional *string* &rarr;  The region in which to obtain the CDN client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *number* &rarr;  ID of the origin group.

- `origins`  *list* &rarr;  List of origin sources in the origin group.
  - `backup` *boolean* &rarr;  Defines whether the origin is a backup, meaning that it will not be used until one of active origins become unavailable.

  - `enabled` *boolean* &rarr;  Enables or disables an origin source in the origin group.

  - `source` *string* &rarr;  IP address or domain name of the origin and the port, if custom port is used.


- `use_next` *boolean* &rarr;  Defines whether to use the next origin from the origin group if origin responds with 4XX or 5XX codes.


