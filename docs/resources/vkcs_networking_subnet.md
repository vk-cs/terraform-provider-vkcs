---
layout: "vkcs"
page_title: "vkcs: networking_subnet"
description: |-
  Manages a subnet resource within VKCS.
---

# vkcs\_networking\_subnet

Manages a subnet resource within VKCS.

## Example Usage

```hcl
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

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the Networking client.
    A Networking client is needed to create a subnet. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    subnet.

* `network_id` - (Required) The UUID of the parent network. Changing this
    creates a new subnet.

* `cidr` - (Optional) CIDR representing IP range for this subnet, based on IP
    version. You can omit this option if you are creating a subnet from a
    subnet pool.

* `prefix_length` - (Optional) The prefix length to use when creating a subnet
    from a subnet pool. The default subnet pool prefix length that was defined
    when creating the subnet pool will be used if not provided. Changing this
    creates a new subnet.

* `ip_version` - (Optional) IP version, either 4 (default) or 6. Changing this creates a
    new subnet.

* `ipv6_address_mode` - (Optional) The IPv6 address mode. Valid values are
  `dhcpv6-stateful`, `dhcpv6-stateless`, or `slaac`.

* `ipv6_ra_mode` - (Optional) The IPv6 Router Advertisement mode. Valid values
  are `dhcpv6-stateful`, `dhcpv6-stateless`, or `slaac`.

* `name` - (Optional) The name of the subnet. Changing this updates the name of
    the existing subnet.

* `description` - (Optional) Human-readable description of the subnet. Changing this
    updates the name of the existing subnet.

* `allocation_pool` - (Optional) A block declaring the start and end range of
    the IP addresses available for use with DHCP in this subnet. Multiple
    `allocation_pool` blocks can be declared, providing the subnet with more
    than one range of IP addresses to use with DHCP. However, each IP range
    must be from the same CIDR that the subnet is part of.
    The `allocation_pool` block is documented below.

* `gateway_ip` - (Optional)  Default gateway used by devices in this subnet.
    Leaving this blank and not setting `no_gateway` will cause a default
    gateway of `.1` to be used. Changing this updates the gateway IP of the
    existing subnet.

* `no_gateway` - (Optional) Do not set a gateway IP on this subnet. Changing
    this removes or adds a default gateway IP of the existing subnet.

* `enable_dhcp` - (Optional) The administrative state of the network.
    Acceptable values are "true" and "false". Changing this value enables or
    disables the DHCP capabilities of the existing subnet. Defaults to true.

* `dns_nameservers` - (Optional) An array of DNS name server names used by hosts
    in this subnet. Changing this updates the DNS name servers for the existing
    subnet.

* `subnetpool_id` - (Optional) The ID of the subnetpool associated with the subnet.

* `value_specs` - (Optional) Map of additional options.

* `tags` - (Optional) A set of string tags for the subnet.

* `sdn` - (Optional) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

The `allocation_pool` block supports:

* `start` - (Required) The starting address.

* `end` - (Required) The ending address.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `network_id` - See Argument Reference above.
* `cidr` - See Argument Reference above.
* `ip_version` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `gateway_ip` - See Argument Reference above.
* `enable_dhcp` - See Argument Reference above.
* `dns_nameservers` - See Argument Reference above.
* `subnetpool_id` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `all_tags` - The collection of ags assigned on the subnet, which have been
  explicitly and implicitly added.
* `sdn` - See Argument Reference above.

## Import

Subnets can be imported using the `id`, e.g.

```
$ terraform import vkcs_networking_subnet.subnet_1 da4faf16-5546-41e4-8330-4d0002b74048
```
