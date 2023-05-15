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
data "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
}
```

## Argument Reference
- `cidr` optional *string* &rarr;  The CIDR of the subnet.

- `description` optional *string* &rarr;  Human-readable description of the subnet.

- `dhcp_enabled` optional *boolean* &rarr;  If the subnet has DHCP enabled.

- `gateway_ip` optional *string* &rarr;  The IP of the subnet's gateway.

- `name` optional *string* &rarr;  The name of the subnet.

- `network_id` optional *string* &rarr;  The ID of the network the subnet belongs to.

- `region` optional *string* &rarr;  The region in which to obtain the Network client. A Network client is needed to retrieve subnet ids. If omitted, the `region` argument of the provider is used.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `subnet_id` optional *string* &rarr;  The ID of the subnet.

- `subnetpool_id` optional *string* &rarr;  The ID of the subnetpool associated with the subnet.

- `tags` optional *set of* *string* &rarr;  The list of subnet tags to filter.

- `tenant_id` optional *string* &rarr;  The owner of the subnet.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `all_tags` *set of* *string* &rarr;  A set of string tags applied on the subnet.

- `allocation_pools` *object* &rarr;  Allocation pools of the subnet.

- `dns_nameservers` *set of* *string* &rarr;  DNS Nameservers of the subnet.

- `enable_dhcp` *boolean* &rarr;  Whether the subnet has DHCP enabled or not.

- `host_routes` *object* &rarr;  Host Routes of the subnet.

- `id` *string* &rarr;  ID of the found subnet.


