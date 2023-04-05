---
layout: "vkcs"
page_title: "vkcs: vkcs_publicdns_zone"
description: |-
  Get information on a public DNS zone.
---

# vkcs_publicdns_zone

Use this data source to get the ID of a VKCS public DNS zone. **New since v.0.2.0**.

## Example Usage

```terraform
data "vkcs_publicdns_zone" "zone" {
  zone = "example.com"
}
```

## Argument Reference
- `admin_email` **String** (*Optional*) The admin email of the zone SOA.

- `expire` **Number** (*Optional*) The expire time of the zone SOA.

- `id` **String** (*Optional*) The UUID of the DNS zone.

- `primary_dns` **String** (*Optional*) The primary DNS of the zone SOA.

- `refresh` **Number** (*Optional*) The refresh time of the zone SOA.

- `region` **String** (*Optional*) The region in which to obtain the V2 Public DNS client. If omitted, the `region` argument of the provider is used.

- `retry` **Number** (*Optional*) The retry time of the zone SOA.

- `serial` **Number** (*Optional*) The serial number of the zone SOA.

- `status` **String** (*Optional*) The status of the zone.

- `ttl` **Number** (*Optional*) The TTL (time to live) of the zone SOA.

- `zone` **String** (*Optional*) The name of the zone.


## Attributes Reference
- `admin_email` **String** See Argument Reference above.

- `expire` **Number** See Argument Reference above.

- `id` **String** The UUID of the DNS zone.

- `primary_dns` **String** See Argument Reference above.

- `refresh` **Number** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `retry` **Number** See Argument Reference above.

- `serial` **Number** See Argument Reference above.

- `status` **String** See Argument Reference above.

- `ttl` **Number** See Argument Reference above.

- `zone` **String** See Argument Reference above.


