---
layout: "vkcs"
page_title: "vkcs: vkcs_networking_network"
description: |-
  Manages a network resource within VKCS.
---

# vkcs_networking_network

Manages a network resource within VKCS.

## Example Usage
```terraform
resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name       = "subnet_1"
  network_id = "${vkcs_networking_network.network_1.id}"
  cidr       = "192.168.199.0/24"
}

resource "vkcs_compute_secgroup" "secgroup_1" {
  name        = "secgroup_1"
  description = "a security group"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "vkcs_networking_port" "port_1" {
  name               = "port_1"
  network_id         = "${vkcs_networking_network.network_1.id}"
  admin_state_up     = "true"
  security_group_ids = ["${vkcs_compute_secgroup.secgroup_1.id}"]

  fixed_ip {
    "subnet_id"  = "${vkcs_networking_subnet.subnet_1.id}"
    "ip_address" = "192.168.199.10"
  }
}

resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["${vkcs_compute_secgroup.secgroup_1.name}"]

  network {
    port = "${vkcs_networking_port.port_1.id}"
  }
}
```
## Argument Reference
- `admin_state_up` **Boolean** (*Optional*) The administrative state of the network. Acceptable values are "true" and "false". Changing this value updates the state of the existing network.

- `description` **String** (*Optional*) Human-readable description of the network. Changing this updates the name of the existing network.

- `name` **String** (*Optional*) The name of the network. Changing this updates the name of the existing network.

- `port_security_enabled` **Boolean** (*Optional*) Whether to explicitly enable or disable port security on the network. Port Security is usually enabled by default, so omitting this argument will usually result in a value of "true". Setting this explicitly to `false` will disable port security. Valid values are `true` and `false`.

- `private_dns_domain` **String** (*Optional*) Private dns domain name

- `region` **String** (*Optional*) The region in which to obtain the Networking client. A Networking client is needed to create a network. If omitted, the `region` argument of the provider is used. Changing this creates a new network.

- `sdn` **String** (*Optional*) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `tags` <strong>Set of </strong>**String** (*Optional*) A set of string tags for the network.

- `value_specs` <strong>Map of </strong>**String** (*Optional*) Map of additional options.

- `vkcs_services_access` **Boolean** (*Optional*) Whether VKCS services access is enabled. This feature should be enabled globally for your project. Access can be enabled for new or existing networks, but cannot be disabled for existing networks. Valid values are `true` and `false`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `all_tags` <strong>Set of </strong>**String** The collection of tags assigned on the network, which have been explicitly and implicitly added.

- `id` **String** ID of the resource.



## Import

Networks can be imported using the `id`, e.g.

```shell
terraform import vkcs_networking_network.network_1 d90ce693-5ccf-4136-a0ed-152ce412b6b9
```
