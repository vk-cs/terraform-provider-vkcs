---
layout: "vkcs"
page_title: "vkcs: vkcs_publicdns_zone"
description: |-
  Manages a public DNS zone.
---

# vkcs_publicdns_zone

Manages a public DNS record resource within VKCS. **New since v.0.2.0**.

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
- `zone` **String** (***Required***) The name of the zone. **Changes this creates a new zone**.

- `admin_email` **String** (*Optional*) The admin email of the zone SOA.

- `expire` **Number** (*Optional*) The expire time of the zone SOA.

- `primary_dns` **String** (*Optional*) The primary DNS of the zone SOA.

- `refresh` **Number** (*Optional*) The refresh time of the zone SOA.

- `region` **String** (*Optional*) The region in which to obtain the V2 Public DNS client. If omitted, the `region` argument of the provider is used. Changing this creates a new zone.

- `retry` **Number** (*Optional*) The retry time of the zone SOA.

- `ttl` **Number** (*Optional*) The TTL (time to live) of the zone SOA.


## Attributes Reference
- `zone` **String** See Argument Reference above.

- `admin_email` **String** See Argument Reference above.

- `expire` **Number** See Argument Reference above.

- `primary_dns` **String** See Argument Reference above.

- `refresh` **Number** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `retry` **Number** See Argument Reference above.

- `ttl` **Number** See Argument Reference above.

- `id` **String** ID of the resource.

- `serial` **Number** The serial number of the zone SOA.

- `status` **String** The status of the zone.



## Import

Public DNS zones can be imported using the `id`, e.g.

```shell
terraform import vkcs_publicdns_zone.zone b758c4e5-ec13-4dfa-8458-b8502625499c
```
