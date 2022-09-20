---
layout: "vkcs"
page_title: "vkcs: vkcs_networking_subnet"
description: |-
  Get information on an VKCS Subnet.
---

# vkcs_networking_subnet

Use this data source to get the ID of an available VKCS subnet.

## Example Usage

```terraform
data "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
}
```

## Argument Reference
- `cidr` **String** (*Optional*) The CIDR of the subnet.

- `description` **String** (*Optional*) Human-readable description of the subnet.

- `dhcp_enabled` **Boolean** (*Optional*) If the subnet has DHCP enabled.

- `gateway_ip` **String** (*Optional*) The IP of the subnet's gateway.

- `name` **String** (*Optional*) The name of the subnet.

- `network_id` **String** (*Optional*) The ID of the network the subnet belongs to.

- `region` **String** (*Optional*) The region in which to obtain the Network client. A Network client is needed to retrieve subnet ids. If omitted, the `region` argument of the provider is used.

- `sdn` **String** (*Optional*) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `subnet_id` **String** (*Optional*) The ID of the subnet.

- `subnetpool_id` **String** (*Optional*) The ID of the subnetpool associated with the subnet.

- `tags` <strong>Set of </strong>**String** (*Optional*) The list of subnet tags to filter.

- `tenant_id` **String** (*Optional*) The owner of the subnet.


## Attributes Reference
- `cidr` **String** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `dhcp_enabled` **Boolean** See Argument Reference above.

- `gateway_ip` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `network_id` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `sdn` **String** See Argument Reference above.

- `subnet_id` **String** See Argument Reference above.

- `subnetpool_id` **String** See Argument Reference above.

- `tags` <strong>Set of </strong>**String** See Argument Reference above.

- `tenant_id` **String** See Argument Reference above.

- `all_tags` <strong>Set of </strong>**String** A set of string tags applied on the subnet.

- `allocation_pools` **Object** Allocation pools of the subnet.

- `dns_nameservers` <strong>Set of </strong>**String** DNS Nameservers of the subnet.

- `enable_dhcp` **Boolean** Whether the subnet has DHCP enabled or not.

- `host_routes` **Object** Host Routes of the subnet.

- `id` **String** ID of the found subnet.


