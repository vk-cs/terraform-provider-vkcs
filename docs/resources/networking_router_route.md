---
subcategory: "Networking"
layout: "vkcs"
page_title: "vkcs: vkcs_networking_router_route"
description: |-
  Creates a routing entry on a VKCS router.
---

# vkcs_networking_router_route

Creates a routing entry on a VKCS router.

## Example Usage
```terraform
resource "vkcs_networking_router_route" "compute_gateway" {
  router_id        = vkcs_networking_router.router.id
  destination_cidr = "10.10.0.0/23"
  next_hop         = vkcs_compute_instance.gateway.access_ip_v4
}
```
## Argument Reference
- `destination_cidr` **required** *string* &rarr;  CIDR block to match on the packet’s destination IP. Changing this creates a new routing entry.

- `next_hop` **required** *string* &rarr;  IP address of the next hop gateway. Changing this creates a new routing entry.

- `router_id` **required** *string* &rarr;  ID of the router this routing entry belongs to. Changing this creates a new routing entry.

- `region` optional *string* &rarr;  The region in which to obtain the networking client. A networking client is needed to configure a routing entry on a router. If omitted, the `region` argument of the provider is used. Changing this creates a new routing entry.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is project's default SDN.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Notes

The `next_hop` IP address must be directly reachable from the router at the ``vkcs_networking_router_route``
resource creation time.  You can ensure that by explicitly specifying a dependency on the ``vkcs_networking_router_interface``
resource that connects the next hop to the router, as in the example above.

## Import

Routing entries can be imported using a combined ID using the following format: ``<router_id>-route-<destination_cidr>-<next_hop>``0

```shell
terraform import vkcs_networking_router_route.router_route_1 686fe248-386c-4f70-9f6c-281607dad079-route-10.0.1.0/24-192.168.199.25
```
