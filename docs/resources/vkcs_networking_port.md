---
layout: "vkcs"
page_title: "vkcs: vkcs_networking_port"
description: |-
  Manages a port resource within VKCS.
---

# vkcs_networking_port

Manages a port resource within VKCS.

## Example Usage
### Simple port
```terraform
resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_port" "port_1" {
  name           = "port_1"
  network_id     = "${vkcs_networking_network.network_1.id}"
  admin_state_up = "true"
}
```

## Argument Reference
- `network_id` **String** (***Required***) The ID of the network to attach the port to. Changing this creates a new port.

- `admin_state_up` **Boolean** (*Optional*) Administrative up/down status for the port (must be `true` or `false` if provided). Changing this updates the `admin_state_up` of an existing port.

- `allowed_address_pairs` (*Optional*) An IP/MAC Address pair of additional IP addresses that can be active on this port. The structure is described below.
  - `ip_address` **String** (***Required***) The additional IP address.

  - `mac_address` **String** (*Optional*) The additional MAC address.

- `description` **String** (*Optional*) Human-readable description of the port. Changing this updates the `description` of an existing port.

- `device_id` **String** (*Optional*) The ID of the device attached to the port. Changing this creates a new port.

- `device_owner` **String** (*Optional*) The device owner of the port. Changing this creates a new port.

- `dns_name` **String** (*Optional*) The port DNS name.

- `extra_dhcp_option` (*Optional*) An extra DHCP option that needs to be configured on the port. The structure is described below. Can be specified multiple times.
  - `name` **String** (***Required***) Name of the DHCP option.

  - `value` **String** (***Required***) Value of the DHCP option.

  - `ip_version` **Number** (*Optional*) IP protocol version. Defaults to 4.

- `fixed_ip` (*Optional*) (Conflicts with `no_fixed_ip`) An array of desired IPs for this port. The structure is described below.
  - `subnet_id` **String** (***Required***) Subnet in which to allocate IP address for this port.

  - `ip_address` **String** (*Optional*) IP address desired in the subnet for this port. If you don't specify `ip_address`, an available IP address from the specified subnet will be allocated to this port. This field will not be populated if it is left blank or omitted. To retrieve the assigned IP address, use the `all_fixed_ips` attribute.

- `mac_address` **String** (*Optional*) Specify a specific MAC address for the port. Changing this creates a new port.

- `name` **String** (*Optional*) A unique name for the port. Changing this updates the `name` of an existing port.

- `no_fixed_ip` **Boolean** (*Optional*) (Conflicts with `fixed_ip`) Create a port with no fixed IP address. This will also remove any fixed IPs previously set on a port. `true` is the only valid value for this argument.

- `no_security_groups` **Boolean** (*Optional*) (Conflicts with `security_group_ids`) If set to `true`, then no security groups are applied to the port. If set to `false` and no `security_group_ids` are specified, then the port will yield to the default behavior of the Networking service, which is to usually apply the "default" security group.

- `port_security_enabled` **Boolean** (*Optional*) Whether to explicitly enable or disable port security on the port. Port Security is usually enabled by default, so omitting argument will usually result in a value of `true`. Setting this explicitly to `false` will disable port security. In order to disable port security, the port must not have any security groups. Valid values are `true` and `false`.

- `region` **String** (*Optional*) The region in which to obtain the Networking client. A Networking client is needed to create a port. If omitted, the `region` argument of the provider is used. Changing this creates a new port.

- `sdn` **String** (*Optional*) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `security_group_ids` <strong>Set of </strong>**String** (*Optional*) (Conflicts with `no_security_groups`) A list of security group IDs to apply to the port. The security groups must be specified by ID and not name (as opposed to how they are configured with the Compute Instance).

- `tags` <strong>Set of </strong>**String** (*Optional*) A set of string tags for the port.

- `value_specs` <strong>Map of </strong>**String** (*Optional*) Map of additional options.


## Attributes Reference
- `network_id` **String** See Argument Reference above.

- `admin_state_up` **Boolean** See Argument Reference above.

- `allowed_address_pairs`  See Argument Reference above.
  - `ip_address` **String** See Argument Reference above.

  - `mac_address` **String** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `device_id` **String** See Argument Reference above.

- `device_owner` **String** See Argument Reference above.

- `dns_name` **String** See Argument Reference above.

- `extra_dhcp_option`  See Argument Reference above.
  - `name` **String** See Argument Reference above.

  - `value` **String** See Argument Reference above.

  - `ip_version` **Number** See Argument Reference above.

- `fixed_ip`  See Argument Reference above.
  - `subnet_id` **String** See Argument Reference above.

  - `ip_address` **String** See Argument Reference above.

- `mac_address` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `no_fixed_ip` **Boolean** See Argument Reference above.

- `no_security_groups` **Boolean** See Argument Reference above.

- `port_security_enabled` **Boolean** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `sdn` **String** See Argument Reference above.

- `security_group_ids` <strong>Set of </strong>**String** See Argument Reference above.

- `tags` <strong>Set of </strong>**String** See Argument Reference above.

- `value_specs` <strong>Map of </strong>**String** See Argument Reference above.

- `all_fixed_ips` **String** The collection of Fixed IP addresses on the port in the order returned by the Network v2 API.

- `all_security_group_ids` <strong>Set of </strong>**String** The collection of Security Group IDs on the port which have been explicitly and implicitly added.

- `all_tags` <strong>Set of </strong>**String** The collection of tags assigned on the port, which have been explicitly and implicitly added.

- `dns_assignment` <strong>Map of </strong>**String** The list of maps representing port DNS assignments.

- `id` **String** ID of the resource.



## Import

Ports can be imported using the `id`, e.g.

```shell
terraform import vkcs_networking_port.port_1 eae26a3e-1c33-4cc1-9c31-0cd729c438a1
```

## Notes

### Ports and Instances

There are some notes to consider when connecting Instances to networks using
Ports. Please see the `vkcs_compute_instance` documentation for further
documentation.
