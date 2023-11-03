---
subcategory: "Network"
layout: "vkcs"
page_title: "vkcs: vkcs_networking_subnet_route"
description: |-
  Creates a routing entry on a VKCS subnet.
---

# vkcs_networking_subnet_route

Creates a routing entry on a VKCS subnet.

## Example Usage
```terraform
resource "vkcs_networking_subnet_route" "subnet-route-to-external-tf-example" {
  subnet_id        = vkcs_networking_subnet.app.id
  destination_cidr = "10.0.1.0/24"
  next_hop         = vkcs_networking_port.persistent_etcd.all_fixed_ips[0]
}
```
## Argument Reference
- `destination_cidr` **required** *string* &rarr;  CIDR block to match on the packet’s destination IP. Changing this creates a new routing entry.

- `next_hop` **required** *string* &rarr;  IP address of the next hop gateway. Changing this creates a new routing entry.

- `subnet_id` **required** *string* &rarr;  ID of the subnet this routing entry belongs to. Changing this creates a new routing entry.

- `region` optional *string* &rarr;  The region in which to obtain the networking client. A networking client is needed to configure a routing entry on a subnet. If omitted, the `region` argument of the provider is used. Changing this creates a new routing entry.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is project's default SDN.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Routing entries can be imported using a combined ID using the following format: ``<subnet_id>-route-<destination_cidr>-<next_hop>``

```shell
terraform import vkcs_networking_subnet_route.subnet_route_1 686fe248-386c-4f70-9f6c-281607dad079-route-10.0.1.0/24-192.168.199.25
```
