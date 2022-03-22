---
layout: "vkcs"
page_title: "VKCS: networking_floatingip"
description: |-
  Manages a V2 floating IP resource within OpenStack Neutron (networking).
---

# vkcs\_networking\_floatingip

Manages a V2 floating IP resource within OpenStack Neutron (networking)
that can be used for load balancers.
These are similar to Nova (compute) floating IP resources,
but only compute floating IPs can be used with compute instances.

## Example Usage

### Simple floating IP allocation

```hcl
resource "vkcs_networking_floatingip" "floatip_1" {
  pool = "public"
}
```

### Floating IP allocation using a list of subnets

If one of the subnets in a list has an exhausted pool, terraform will try the
next subnet ID from the list.

```hcl
data "vkcs_networking_network" "ext_network" {
  name = "public"
}

data "openstack_networking_subnet_ids_v2" "ext_subnets" {
  network_id = data.vkcs_networking_network.ext_network.id
}

resource "vkcs_networking_floatingip" "floatip_1" {
  pool       = data.vkcs_networking_network.ext_network.name
  subnet_ids = data.openstack_networking_subnet_ids_v2.ext_subnets.ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create a floating IP that can be used with
  another networking resource, such as a load balancer. If omitted, the
  `region` argument of the provider is used. Changing this creates a new
  floating IP (which may or may not have a different address).

* `description` - (Optional) Human-readable description for the floating IP.

* `pool` - (Required) The name of the pool from which to obtain the floating
  IP. Changing this creates a new floating IP.

* `port_id` - (Optional) ID of an existing port with at least one IP address to
  associate with this floating IP.

* `fixed_ip` - Fixed IP of the port to associate with this floating IP. Required if
  the port has multiple fixed IPs.

* `subnet_id` - (Optional) The subnet ID of the floating IP pool. Specify this if
  the floating IP network has multiple subnets.

* `subnet_ids` - (Optional) A list of external subnet IDs to try over each to
  allocate a floating IP address. If a subnet ID in a list has exhausted
  floating IP pool, the next subnet ID will be tried. This argument is used only
  during the resource creation. Conflicts with a `subnet_id` argument.

* `value_specs` - (Optional) Map of additional options.

* `tags` - (Optional) A set of string tags for the floating IP.

* `sdn` - (Optional) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `description` - See Argument Reference above.
* `pool` - See Argument Reference above.
* `address` - The actual floating IP address itself.
* `port_id` - ID of associated port.
* `tenant_id` - the ID of the tenant in which to create the floating IP.
* `fixed_ip` - The fixed IP which the floating IP maps to.
* `tags` - See Argument Reference above.
* `all_tags` - The collection of tags assigned on the floating IP, which have
  been explicitly and implicitly added.
* `sdn` - See Argument Reference above.

## Import

Floating IPs can be imported using the `id`, e.g.

```
$ terraform import vkcs_networking_floatingip.floatip_1 2c7f39f3-702b-48d1-940c-b50384177ee1
```
