---
subcategory: "Network"
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
- `pool` **required** *string* &rarr;  The name of the pool from which to obtain the floating IP. Changing this creates a new floating IP.

- `address` optional *string* &rarr;  The actual floating IP address itself.

- `description` optional *string* &rarr;  Human-readable description for the floating IP.

- `fixed_ip` optional *string* &rarr;  Fixed IP of the port to associate with this floating IP. Required if the port has multiple fixed IPs.

- `port_id` optional *string* &rarr;  ID of an existing port with at least one IP address to associate with this floating IP.

- `region` optional *string* &rarr;  The region in which to obtain the Networking client. A Networking client is needed to create a floating IP that can be used with another networking resource, such as a load balancer. If omitted, the `region` argument of the provider is used. Changing this creates a new floating IP (which may or may not have a different address).

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `subnet_id` optional *string* &rarr;  The subnet ID of the floating IP pool. Specify this if the floating IP network has multiple subnets.

- `subnet_ids` optional *string* &rarr;  A list of external subnet IDs to try over each to allocate a floating IP address. If a subnet ID in a list has exhausted floating IP pool, the next subnet ID will be tried. This argument is used only during the resource creation. Conflicts with a `subnet_id` argument.

- `value_specs` optional *map of* *string* &rarr;  Map of additional options.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Floating IPs can be imported using the `id`, e.g.

```shell
terraform import vkcs_networking_floatingip.floatip_1 2c7f39f3-702b-48d1-940c-b50384177ee1
```
