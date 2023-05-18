---
subcategory: "DNS"
layout: "vkcs"
page_title: "vkcs: vkcs_publicdns_zone"
description: |-
  Manages a public DNS zone.
---

# vkcs_publicdns_zone

Manages a public DNS record resource within VKCS. **New since v0.2.0**.

## Example Usage
```terraform
resource "vkcs_publicdns_zone" "zone" {
  zone = local.zone_name
  primary_dns = "ns1.mcs.mail.ru"
  admin_email = "admin@example.com"
  expire = 3600000
}
```
## Argument Reference
- `zone` **required** *string* &rarr;  The name of the zone. **Changes this creates a new zone**.

- `admin_email` optional *string* &rarr;  The admin email of the zone SOA.

- `expire` optional *number* &rarr;  The expire time of the zone SOA.

- `primary_dns` optional *string* &rarr;  The primary DNS of the zone SOA.

- `refresh` optional *number* &rarr;  The refresh time of the zone SOA.

- `region` optional *string* &rarr;  The region in which to obtain the V2 Public DNS client. If omitted, the `region` argument of the provider is used. Changing this creates a new zone.

- `retry` optional *number* &rarr;  The retry time of the zone SOA.

- `ttl` optional *number* &rarr;  The TTL (time to live) of the zone SOA.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `serial` *number* &rarr;  The serial number of the zone SOA.

- `status` *string* &rarr;  The status of the zone.



## Import

Public DNS zones can be imported using the `id`, e.g.

```shell
terraform import vkcs_publicdns_zone.zone b758c4e5-ec13-4dfa-8458-b8502625499c
```
