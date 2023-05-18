---
subcategory: "DNS"
layout: "vkcs"
page_title: "vkcs: vkcs_publicdns_zone"
description: |-
  Get information on a public DNS zone.
---

# vkcs_publicdns_zone

Use this data source to get the ID of a VKCS public DNS zone. **New since v0.2.0**.

## Example Usage

```terraform
data "vkcs_publicdns_zone" "zone" {
  zone = "example.com"
}
```

## Argument Reference
- `admin_email` optional *string* &rarr;  The admin email of the zone SOA.

- `expire` optional *number* &rarr;  The expire time of the zone SOA.

- `id` optional *string* &rarr;  The UUID of the DNS zone.

- `primary_dns` optional *string* &rarr;  The primary DNS of the zone SOA.

- `refresh` optional *number* &rarr;  The refresh time of the zone SOA.

- `region` optional *string* &rarr;  The region in which to obtain the V2 Public DNS client. If omitted, the `region` argument of the provider is used.

- `retry` optional *number* &rarr;  The retry time of the zone SOA.

- `serial` optional *number* &rarr;  The serial number of the zone SOA.

- `status` optional *string* &rarr;  The status of the zone.

- `ttl` optional *number* &rarr;  The TTL (time to live) of the zone SOA.

- `zone` optional *string* &rarr;  The name of the zone.


## Attributes Reference
No additional attributes are exported.

