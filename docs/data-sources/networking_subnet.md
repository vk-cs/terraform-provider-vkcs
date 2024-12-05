---
subcategory: "Network"
layout: "vkcs"
page_title: "vkcs: vkcs_networking_subnet"
description: |-
  Get information on an VKCS Subnet.
---

# vkcs_networking_subnet

Use this data source to get the ID of an available VKCS subnet.

## Example Usage

```terraform
data "vkcs_networking_subnet" "subnet_one_of_internal" {
  cidr       = "192.168.199.0/24"
  network_id = vkcs_networking_network.app.id
  # This is unnecessary in real life.
  # This is required here to let the example work with subnet resource example. 
  depends_on = [vkcs_networking_subnet.app]
}
```

## Argument Reference
- `cidr` optional *string* &rarr;  The CIDR of the subnet.

- `description` optional *string* &rarr;  Human-readable description of the subnet.

- `dhcp_enabled` optional *boolean* &rarr;  If the subnet has DHCP enabled.

- `gateway_ip` optional *string* &rarr;  The IP of the subnet's gateway.

- `id` optional *string* &rarr;  The ID of the subnet.

- `name` optional *string* &rarr;  The name of the subnet.

- `network_id` optional *string* &rarr;  The ID of the network the subnet belongs to.

- `region` optional *string* &rarr;  The region in which to obtain the Network client. A Network client is needed to retrieve subnet ids. If omitted, the `region` argument of the provider is used.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is project's default SDN.

- `subnet_id` optional deprecated *string* &rarr;  The ID of the subnet. **Deprecated** This argument is deprecated, please, use the `id` attribute instead.

- `subnetpool_id` optional *string* &rarr;  The ID of the subnetpool associated with the subnet.

- `tags` optional *set of* *string* &rarr;  The list of subnet tags to filter.

- `tenant_id` optional *string* &rarr;  The owner of the subnet.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `all_tags` *set of* *string* &rarr;  A set of string tags applied on the subnet.

- `allocation_pools`  *list* &rarr;  Allocation pools of the subnet.
  - `end` *string* &rarr;  The ending address.

  - `start` *string* &rarr;  The starting address.


- `dns_nameservers` *set of* *string* &rarr;  DNS Nameservers of the subnet.

- `enable_dhcp` *boolean* &rarr;  Whether the subnet has DHCP enabled or not.

- `host_routes`  *list* &rarr;  Host Routes of the subnet.
  - `destination_cidr` *string*

  - `next_hop` *string*



