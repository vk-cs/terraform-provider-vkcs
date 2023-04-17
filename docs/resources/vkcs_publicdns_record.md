---
layout: "vkcs"
page_title: "vkcs: vkcs_publicdns_record"
description: |-
  Manages a public DNS record.
---

# vkcs_publicdns_record

Manages a public DNS zone record resource within VKCS. **New since v.0.2.0**.<br>
**Note:** Although some arguments are marked as optional, it is actually required to set values for them depending on record \"type\". Use this map to get information about which arguments you have to set:

| Record type | Required arguments |
| ----------- | ------------------ |
| A | ip |
| AAAA | ip |
| CNAME | name, content |
| MX | priority, content |
| NS | content |
| SRV | service, proto, priority, weight, host, port |
| TXT | content |




## Example Usage
### Single record
```terraform
resource "vkcs_publicdns_record" "srv" {
  zone_id = vkcs_publicdns_zone.zone.id
  type = "SRV"
  service = "_sip"
  proto = "_udp"
  priority = 10
  weight = 5
  host = "siptarget.com"
  port = 5060
  ttl = 60
}
```

### Multiple A records
```terraform
locals {
  google_public_dns_ips = tomap({
    "ip_1" = "8.8.8.8"
    "ip_2" = "8.8.4.4"
  })
}

resource "vkcs_publicdns_record" "a" {
  for_each = local.google_public_dns_ips
  zone_id = vkcs_publicdns_zone.zone.id
  type = "A"
  name = "google-dns-servers"
  ip = each.value
  ttl = 60
}
```

## Argument Reference
- `type` **String** (***Required***) The type of the record. Must be one of following: "A", "AAAA", "CNAME", "MX", "NS", "SRV", "TXT".

- `zone_id` **String** (***Required***) The ID of the zone to attach the record to.

- `content` **String** (*Optional*) The content of the record.

- `host` **String** (*Optional*) The domain name of the target host.

- `ip` **String** (*Optional*) The IP address of the record. It should be IPv4 for record of type "A" and IPv6 for record of type "AAAA".

- `name` **String** (*Optional*) The name of the record.

- `port` **Number** (*Optional*) The port on the target host of the service.

- `priority` **Number** (*Optional*) The priority of the record's server.

- `proto` **String** (*Optional*) The name of the desired protocol.

- `region` **String** (*Optional*) The region in which to obtain the V2 Public DNS client. If omitted, the `region` argument of the provider is used. Changing this creates a new record.

- `service` **String** (*Optional*) The name of the desired service.

- `ttl` **Number** (*Optional*) The time to live of the record.

- `weight` **Number** (*Optional*) The relative weight of the record's server.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `full_name` **String** The full name of the SRV record.

- `id` **String** ID of the resource.


## Import

Public DNS records can be imported using the `id` in the form `<zone-id>/<record-type>/<record-id>`, e.g.

```shell
terraform import vkcs_publicdns_record.record 7582c61b-99b7-4730-a74f-7062fbadb94c/a/96b11adf-2627-4a06-bceb-a7f3b61b709e
```
