---
subcategory: "Network"
layout: "vkcs"
page_title: "vkcs: vkcs_networking_router_interface"
description: |-
  Manages a router interface resource within VKCS.
---

# vkcs_networking_router_interface

Manages a router interface resource within VKCS.

## Example Usage
```terraform
resource "vkcs_networking_router_interface" "app" {
  router_id = vkcs_networking_router.router.id
  subnet_id = vkcs_networking_subnet.app.id
}
```

## Argument Reference
- `router_id` **required** *string* &rarr;  ID of the router this interface belongs to. Changing this creates a new router interface.

- `port_id` optional *string* &rarr;  ID of the port this interface connects to. Changing this creates a new router interface.

- `region` optional *string* &rarr;  The region in which to obtain the networking client. A networking client is needed to create a router. If omitted, the `region` argument of the provider is used. Changing this creates a new router interface.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is project's default SDN.

- `subnet_id` optional *string* &rarr;  ID of the subnet this interface connects to. Changing this creates a new router interface.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Router Interfaces can be imported using the port `id`, e.g.

```shell
openstack port list --router <router name or id>
terraform import vkcs_networking_router_interface.int_1 <port id from above output>
```
