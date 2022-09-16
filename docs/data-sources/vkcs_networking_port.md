---
layout: "vkcs"
page_title: "vkcs: vkcs_networking_port"
description: |-
  Get information of an VKCS Port.
---

# vkcs_networking_port

Use this data source to get the ID of an available VKCS port.

## Example Usage

```terraform
data "vkcs_networking_port" "port_1" {
  name = "port_1"
}
```

## Argument Reference
- `admin_state_up` **Boolean** (*Optional*) The administrative state of the port.

- `description` **String** (*Optional*) Human-readable description of the port.

- `device_id` **String** (*Optional*) The ID of the device the port belongs to.

- `device_owner` **String** (*Optional*) The device owner of the port.

- `dns_name` **String** (*Optional*) The port DNS name to filter.

- `fixed_ip` **String** (*Optional*) The port IP address filter.

- `mac_address` **String** (*Optional*) The MAC address of the port.

- `name` **String** (*Optional*) The name of the port.

- `network_id` **String** (*Optional*) The ID of the network the port belongs to.

- `port_id` **String** (*Optional*) The ID of the port.

- `project_id` **String** (*Optional*) The project_id of the owner of the port.

- `region` **String** (*Optional*) The region in which to obtain the Network client. A Network client is needed to retrieve port ids. If omitted, the `region` argument of the provider is used.

- `sdn` **String** (*Optional*) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `security_group_ids` <strong>Set of </strong>**String** (*Optional*) The list of port security group IDs to filter.

- `status` **String** (*Optional*) The status of the port.

- `tags` <strong>Set of </strong>**String** (*Optional*) The list of port tags to filter.

- `tenant_id` **String** (*Optional*) The tenant_id of the owner of the port.


## Attributes Reference
- `admin_state_up` **Boolean** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `device_id` **String** See Argument Reference above.

- `device_owner` **String** See Argument Reference above.

- `dns_name` **String** See Argument Reference above.

- `fixed_ip` **String** See Argument Reference above.

- `mac_address` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `network_id` **String** See Argument Reference above.

- `port_id` **String** See Argument Reference above.

- `project_id` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `sdn` **String** See Argument Reference above.

- `security_group_ids` <strong>Set of </strong>**String** See Argument Reference above.

- `status` **String** See Argument Reference above.

- `tags` <strong>Set of </strong>**String** See Argument Reference above.

- `tenant_id` **String** See Argument Reference above.

- `all_fixed_ips` **String** The collection of Fixed IP addresses on the port in the order returned by the Network v2 API.

- `all_security_group_ids` <strong>Set of </strong>**String** The set of security group IDs applied on the port.

- `all_tags` <strong>Set of </strong>**String** The set of string tags applied on the port.

- `allowed_address_pairs` <strong>Set of </strong>**Object** An IP/MAC Address pair of additional IP addresses that can be active on this port. The structure is described below.

- `dns_assignment` <strong>Map of </strong>**String** The list of maps representing port DNS assignments.

- `extra_dhcp_option` **Object** An extra DHCP option configured on the port. The structure is described below.

- `id` **String** ID of the found port.


