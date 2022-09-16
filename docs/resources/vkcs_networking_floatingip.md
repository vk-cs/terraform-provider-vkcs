---
layout: "vkcs"
page_title: "vkcs: vkcs_networking_floatingip"
description: |-
  Manages a floating IP resource within VKCS.
---

# vkcs_networking_floatingip

Manages a floating IP resource within VKCS that can be used for load balancers.

## Example Usage
### Simple floating IP allocation
```terraform
resource "vkcs_networking_floatingip" "floatip_1" {
  pool = "public"
}
```

### Floating IP allocation using a list of subnets
If one of the subnets in a list has an exhausted pool, terraform will try the
next subnet ID from the list.

```terraform
data "vkcs_networking_network" "ext_network" {
  name = "public"
}

resource "vkcs_networking_floatingip" "floatip_1" {
  pool       = data.vkcs_networking_network.ext_network.name
  subnet_ids = [<subnet1_id>, <subnet2_id>, <subnet3_id>]
}
```

## Argument Reference
- `pool` **String** (***Required***) The name of the pool from which to obtain the floating IP. Changing this creates a new floating IP.

- `address` **String** (*Optional*) The actual floating IP address itself.

- `description` **String** (*Optional*) Human-readable description for the floating IP.

- `fixed_ip` **String** (*Optional*) Fixed IP of the port to associate with this floating IP. Required if the port has multiple fixed IPs.

- `port_id` **String** (*Optional*) ID of an existing port with at least one IP address to associate with this floating IP.

- `region` **String** (*Optional*) The region in which to obtain the Networking client. A Networking client is needed to create a floating IP that can be used with another networking resource, such as a load balancer. If omitted, the `region` argument of the provider is used. Changing this creates a new floating IP (which may or may not have a different address).

- `sdn` **String** (*Optional*) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `subnet_id` **String** (*Optional*) The subnet ID of the floating IP pool. Specify this if the floating IP network has multiple subnets.

- `subnet_ids` **String** (*Optional*) A list of external subnet IDs to try over each to allocate a floating IP address. If a subnet ID in a list has exhausted floating IP pool, the next subnet ID will be tried. This argument is used only during the resource creation. Conflicts with a `subnet_id` argument.

- `value_specs` <strong>Map of </strong>**String** (*Optional*) Map of additional options.


## Attributes Reference
- `pool` **String** See Argument Reference above.

- `address` **String** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `fixed_ip` **String** See Argument Reference above.

- `port_id` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `sdn` **String** See Argument Reference above.

- `subnet_id` **String** See Argument Reference above.

- `subnet_ids` **String** See Argument Reference above.

- `value_specs` <strong>Map of </strong>**String** See Argument Reference above.

- `id` **String** ID of the resource.



## Import

Floating IPs can be imported using the `id`, e.g.

```shell
terraform import vkcs_networking_floatingip.floatip_1 2c7f39f3-702b-48d1-940c-b50384177ee1
```
