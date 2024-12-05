---
subcategory: "Network"
layout: "vkcs"
page_title: "vkcs: vkcs_networking_port"
description: |-
  Get information of an VKCS Port.
---

# vkcs_networking_port

Use this data source to get the ID of an available VKCS port.

## Example Usage

```terraform
data "vkcs_networking_port" "persistent_etcd" {
  tags       = ["tf-example", "etcd"]
  network_id = vkcs_networking_network.db.id
  # This is unnecessary in real life.
  # This is required here to let the example work with port resource example. 
  depends_on = [vkcs_networking_port.persistent_etcd]
}
```

## Argument Reference
- `admin_state_up` optional *boolean* &rarr;  The administrative state of the port.

- `description` optional *string* &rarr;  Human-readable description of the port.

- `device_id` optional *string* &rarr;  The ID of the device the port belongs to.

- `device_owner` optional *string* &rarr;  The device owner of the port.

- `dns_name` optional *string* &rarr;  The port DNS name to filter.

- `fixed_ip` optional *string* &rarr;  The port IP address filter.

- `id` optional *string* &rarr;  The ID of the port.

- `mac_address` optional *string* &rarr;  The MAC address of the port.

- `name` optional *string* &rarr;  The name of the port.

- `network_id` optional *string* &rarr;  The ID of the network the port belongs to.

- `port_id` optional deprecated *string* &rarr;  The ID of the port. **Deprecated** This argument is deprecated, please, use the `id` attribute instead.

- `project_id` optional *string* &rarr;  The project_id of the owner of the port.

- `region` optional *string* &rarr;  The region in which to obtain the Network client. A Network client is needed to retrieve port ids. If omitted, the `region` argument of the provider is used.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is project's default SDN.

- `security_group_ids` optional *set of* *string* &rarr;  The list of port security group IDs to filter.

- `status` optional *string* &rarr;  The status of the port.

- `tags` optional *set of* *string* &rarr;  The list of port tags to filter.

- `tenant_id` optional *string* &rarr;  The tenant_id of the owner of the port.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `all_fixed_ips` *string* &rarr;  The collection of Fixed IP addresses on the port in the order returned by the Network v2 API.

- `all_security_group_ids` *set of* *string* &rarr;  The set of security group IDs applied on the port.

- `all_tags` *set of* *string* &rarr;  The set of string tags applied on the port.

- `allowed_address_pairs`  *set* &rarr;  An IP/MAC Address pair of additional IP addresses that can be active on this port. The structure is described below.
  - `ip_address` *string* &rarr;  The additional IP address.

  - `mac_address` *string* &rarr;  The additional MAC address.


- `dns_assignment` *map of* *string* &rarr;  The list of maps representing port DNS assignments.

- `extra_dhcp_option`  *list* &rarr;  An extra DHCP option configured on the port. The structure is described below.
  - `name` *string* &rarr;  Name of the DHCP option.

  - `value` *string* &rarr;  Value of the DHCP option.



