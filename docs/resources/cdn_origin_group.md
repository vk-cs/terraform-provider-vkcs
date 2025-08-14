---
subcategory: "CDN"
layout: "vkcs"
page_title: "vkcs: vkcs_cdn_origin_group"
description: |-
  Manages a CDN origin group within VKCS.
---

# vkcs_cdn_origin_group



## Example Usage
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

## Argument Reference
- `name` **required** *string* &rarr;  Name of the origin group.

- `origins`  *list* &rarr;  List of origin sources in the origin group.
    - `source` **required** *string* &rarr;  IP address or domain name of the origin and the port, if custom port is used.

    - `backup` optional *boolean* &rarr;  Defines whether the origin is a backup, meaning that it will not be used until one of active origins become unavailable. Defaults to false.

    - `enabled` optional *boolean* &rarr;  Enables or disables an origin source in the origin group. Enabled by default.


- `region` optional *string* &rarr;  The region in which to obtain the CDN client. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.

- `use_next` optional *boolean* &rarr;  Defines whether to use the next origin from the origin group if origin responds with 4XX or 5XX codes. Defaults to false.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *number* &rarr;  ID of the origin group.



## Import

An origin group can be imported using the `id`, e.g.
```shell
terraform import vkcs_cdn_resource.resource <resource_id>
```
