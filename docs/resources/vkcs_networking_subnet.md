---
layout: "vkcs"
page_title: "vkcs: vkcs_networking_subnet"
description: |-
  Manages a subnet resource within VKCS.
---

# vkcs_networking_subnet

Manages a subnet resource within VKCS.

## Example Usage
```terraform
resource "vkcs_networking_network" "network_1" {
  name           = "tf_test_network"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  network_id = "${vkcs_networking_network.network_1.id}"
  cidr       = "192.168.199.0/24"
}
```
## Argument Reference
- `network_id` **String** (***Required***) The UUID of the parent network. Changing this creates a new subnet.

- `allocation_pool` (*Optional*) A block declaring the start and end range of the IP addresses available for use with DHCP in this subnet. Multiple `allocation_pool` blocks can be declared, providing the subnet with more than one range of IP addresses to use with DHCP. However, each IP range must be from the same CIDR that the subnet is part of. The `allocation_pool` block is documented below.
  - `end` **String** (***Required***) The ending address.

  - `start` **String** (***Required***) The starting address.

- `cidr` **String** (*Optional*) CIDR representing IP range for this subnet, based on IP version. You can omit this option if you are creating a subnet from a subnet pool.

- `description` **String** (*Optional*) Human-readable description of the subnet. Changing this updates the name of the existing subnet.

- `dns_nameservers` **String** (*Optional*) An array of DNS name server names used by hosts in this subnet. Changing this updates the DNS name servers for the existing subnet.

- `enable_dhcp` **Boolean** (*Optional*) The administrative state of the network. Acceptable values are "true" and "false". Changing this value enables or disables the DHCP capabilities of the existing subnet. Defaults to true.

- `gateway_ip` **String** (*Optional*) Default gateway used by devices in this subnet. Leaving this blank and not setting `no_gateway` will cause a default gateway of `.1` to be used. Changing this updates the gateway IP of the existing subnet.

- `ip_version` **Number** (*Optional*) IP version, either 4 (default) or 6. Changing this creates a new subnet.

- `name` **String** (*Optional*) The name of the subnet. Changing this updates the name of the existing subnet.

- `no_gateway` **Boolean** (*Optional*) Do not set a gateway IP on this subnet. Changing this removes or adds a default gateway IP of the existing subnet.

- `prefix_length` **Number** (*Optional*) The prefix length to use when creating a subnet from a subnet pool. The default subnet pool prefix length that was defined when creating the subnet pool will be used if not provided. Changing this creates a new subnet.

- `region` **String** (*Optional*) The region in which to obtain the Networking client. A Networking client is needed to create a subnet. If omitted, the `region` argument of the provider is used. Changing this creates a new subnet.

- `sdn` **String** (*Optional*) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `subnetpool_id` **String** (*Optional*) The ID of the subnetpool associated with the subnet.

- `tags` <strong>Set of </strong>**String** (*Optional*) A set of string tags for the subnet.

- `value_specs` <strong>Map of </strong>**String** (*Optional*) Map of additional options.


## Attributes Reference
- `network_id` **String** See Argument Reference above.

- `allocation_pool`  See Argument Reference above.
  - `end` **String** See Argument Reference above.

  - `start` **String** See Argument Reference above.

- `cidr` **String** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `dns_nameservers` **String** See Argument Reference above.

- `enable_dhcp` **Boolean** See Argument Reference above.

- `gateway_ip` **String** See Argument Reference above.

- `ip_version` **Number** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `no_gateway` **Boolean** See Argument Reference above.

- `prefix_length` **Number** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `sdn` **String** See Argument Reference above.

- `subnetpool_id` **String** See Argument Reference above.

- `tags` <strong>Set of </strong>**String** See Argument Reference above.

- `value_specs` <strong>Map of </strong>**String** See Argument Reference above.

- `all_tags` <strong>Set of </strong>**String** The collection of ags assigned on the subnet, which have been explicitly and implicitly added.

- `id` **String** ID of the resource.



## Import

Subnets can be imported using the `id`, e.g.

```shell
terraform import vkcs_networking_subnet.subnet_1 da4faf16-5546-41e4-8330-4d0002b74048
```
