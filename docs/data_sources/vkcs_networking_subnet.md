---
layout: "vkcs"
page_title: "vkcs: networking_subnet"
description: |-
  Get information on an VKCS Subnet.
---

# vkcs\_networking\_subnet

Use this data source to get the ID of an available VKCS subnet.

## Example Usage

```hcl
data "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the Network client.
  A Network client is needed to retrieve subnet ids. If omitted, the
  `region` argument of the provider is used.

* `name` - (Optional) The name of the subnet.

* `description` - (Optional) Human-readable description of the subnet.

* `dhcp_enabled` - (Optional) If the subnet has DHCP enabled.

* `network_id` - (Optional) The ID of the network the subnet belongs to.

* `tenant_id` - (Optional) The owner of the subnet.

* `ip_version` - (Optional) The IP version of the subnet (either 4 or 6).

* `gateway_ip` - (Optional) The IP of the subnet's gateway.

* `cidr` - (Optional) The CIDR of the subnet.

* `subnet_id` - (Optional) The ID of the subnet.

* `ipv6_address_mode` - (Optional) The IPv6 address mode. Valid values are
  `dhcpv6-stateful`, `dhcpv6-stateless`, or `slaac`.

* `ipv6_ra_mode` - (Optional) The IPv6 Router Advertisement mode. Valid values
  are `dhcpv6-stateful`, `dhcpv6-stateless`, or `slaac`.

* `subnetpool_id` - (Optional) The ID of the subnetpool associated with the subnet.

* `tags` - (Optional) The list of subnet tags to filter.

* `sdn` - (Optional) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

## Attributes Reference

`id` is set to the ID of the found subnet. In addition, the following attributes
are exported:

* `allocation_pools` - Allocation pools of the subnet.
* `enable_dhcp` - Whether the subnet has DHCP enabled or not.
* `dns_nameservers` - DNS Nameservers of the subnet.
* `host_routes` - Host Routes of the subnet.
* `region` - See Argument Reference above.
* `all_tags` - A set of string tags applied on the subnet.
* `sdn` - See Argument Reference above.
