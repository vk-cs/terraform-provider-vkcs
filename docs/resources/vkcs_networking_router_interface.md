---
layout: "vkcs"
page_title: "vkcs: networking_router_interface"
description: |-
  Manages a router interface resource within VKCS.
---

# vkcs\_networking\_router\_interface

Manages a router interface resource within VKCS.

## Example Usage

```hcl
resource "vkcs_networking_network" "network_1" {
  name           = "tf_test_network"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  network_id = "${vkcs_networking_network.network_1.id}"
  cidr       = "192.168.199.0/24"
  ip_version = 4
}

resource "vkcs_networking_router" "router_1" {
  name                = "my_router"
  external_network_id = "f67f0d72-0ddf-11e4-9d95-e1f29f417e2f"
}

resource "vkcs_networking_router_interface" "router_interface_1" {
  router_id = "${vkcs_networking_router.router_1.id}"
  subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the networking client.
    A networking client is needed to create a router. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    router interface.

* `router_id` - (Required) ID of the router this interface belongs to. Changing
    this creates a new router interface.

* `subnet_id` - ID of the subnet this interface connects to. Changing
    this creates a new router interface.

* `port_id` - ID of the port this interface connects to. Changing
    this creates a new router interface.

* `sdn` - (Optional) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `router_id` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `port_id` - See Argument Reference above.
* `sdn` - See Argument Reference above.

## Import

Router Interfaces can be imported using the port `id`, e.g.

```
$ openstack port list --router <router name or id>
$ terraform import vkcs_networking_router_interface.int_1 <port id from above output>
```
