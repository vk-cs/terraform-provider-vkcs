---
subcategory: "Network"
layout: "vkcs"
page_title: "vkcs: vkcs_networking_floatingip_associate"
description: |-
  Associates a Floating IP to a Port
---

# vkcs_networking_floatingip_associate

Associates a floating IP to a port. This can be done only if port is assigned to router connected to external network. This is useful for situations where you have a pre-allocated floating IP or are unable to use the `vkcs_networking_floatingip` resource to create a floating IP.

## Example Usage
```terraform
resource "vkcs_networking_floatingip_associate" "floatingip_associate" {
  floating_ip = vkcs_networking_floatingip.base_fip.address
  port_id     = vkcs_networking_port.persistent_etcd.id
  # Ensure the router interface is up
  depends_on = [vkcs_networking_router_interface.db]
}
```
## Argument Reference
- `floating_ip` **required** *string* &rarr;  IP Address of an existing floating IP.

- `port_id` **required** *string* &rarr;  ID of an existing port with at least one IP address to associate with this floating IP.

- `fixed_ip` optional *string* &rarr;  One of the port's IP addresses.

- `region` optional *string* &rarr;  The region in which to obtain the Networking client. A Networking client is needed to create a floating IP that can be used with another networking resource, such as a load balancer. If omitted, the `region` argument of the provider is used. Changing this creates a new floating IP (which may or may not have a different address).

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is project's default SDN.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Floating IP associations can be imported using the `id` of the floating IP, e.g.

```shell
terraform import vkcs_networking_floatingip_associate.fip 2c7f39f3-702b-48d1-940c-b50384177ee1
```
