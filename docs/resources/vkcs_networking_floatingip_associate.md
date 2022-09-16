---
layout: "vkcs"
page_title: "vkcs: vkcs_networking_floatingip_associate"
description: |-
  Associates a Floating IP to a Port
---

# vkcs_networking_floatingip_associate

Associates a floating IP to a port. This is useful for situations where you have a pre-allocated floating IP or are unable to use the `vkcs_networking_floatingip` resource to create a floating IP.

## Example Usage
```terraform
resource "vkcs_networking_port" "port_1" {
  network_id = "a5bbd213-e1d3-49b6-aed1-9df60ea94b9a"
}

resource "vkcs_networking_floatingip_associate" "fip_1" {
  floating_ip = "1.2.3.4"
  port_id     = "${vkcs_networking_port.port_1.id}"
}
```
## Argument Reference
- `floating_ip` **String** (***Required***) IP Address of an existing floating IP.

- `port_id` **String** (***Required***) ID of an existing port with at least one IP address to associate with this floating IP.

- `fixed_ip` **String** (*Optional*) One of the port's IP addresses.

- `region` **String** (*Optional*) The region in which to obtain the Networking client. A Networking client is needed to create a floating IP that can be used with another networking resource, such as a load balancer. If omitted, the `region` argument of the provider is used. Changing this creates a new floating IP (which may or may not have a different address).

- `sdn` **String** (*Optional*) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".


## Attributes Reference
- `floating_ip` **String** See Argument Reference above.

- `port_id` **String** See Argument Reference above.

- `fixed_ip` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `sdn` **String** See Argument Reference above.

- `id` **String** ID of the resource.



## Import

Floating IP associations can be imported using the `id` of the floating IP, e.g.

```shell
terraform import vkcs_networking_floatingip_associate.fip 2c7f39f3-702b-48d1-940c-b50384177ee1
```
