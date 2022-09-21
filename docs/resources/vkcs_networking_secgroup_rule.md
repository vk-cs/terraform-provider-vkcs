---
layout: "vkcs"
page_title: "vkcs: vkcs_networking_secgroup_rule"
description: |-
  Manages a security group rule resource within VKCS.
---

# vkcs_networking_secgroup_rule

Manages a security group rule resource within VKCS.

## Example Usage
```terraform
resource "vkcs_networking_secgroup" "secgroup_1" {
  name        = "secgroup_1"
  description = "My security group"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}
```

## Argument Reference
- `direction` **String** (***Required***) The direction of the rule, valid values are __ingress__ or __egress__. Changing this creates a new security group rule.

- `ethertype` **String** (***Required***) The layer 3 protocol type, valid values are __IPv4__ or __IPv6__. Changing this creates a new security group rule.

- `security_group_id` **String** (***Required***) The security group id the rule should belong to, the value needs to be an ID of a security group in the same tenant. Changing this creates a new security group rule.

- `description` **String** (*Optional*) A description of the rule. Changing this creates a new security group rule.

- `port_range_max` **Number** (*Optional*) The higher part of the allowed port range, valid integer value needs to be between 1 and 65535. Changing this creates a new security group rule.

- `port_range_min` **Number** (*Optional*) The lower part of the allowed port range, valid integer value needs to be between 1 and 65535. Changing this creates a new security group rule.

- `protocol` **String** (*Optional*) The layer 4 protocol type, valid values are following. Changing this creates a new security group rule. This is required if you want to specify a port range.
  * __tcp__
  * __udp__
  * __icmp__
  * __ah__
  * __dccp__
  * __egp__
  * __esp__
  * __gre__
  * __igmp__
  * __ospf__
  * __pgm__
  * __rsvp__
  * __sctp__
  * __udplite__
  * __vrrp__

- `region` **String** (*Optional*) The region in which to obtain the networking client. A networking client is needed to create a port. If omitted, the `region` argument of the provider is used. Changing this creates a new security group rule.

- `remote_group_id` **String** (*Optional*) The remote group id, the value needs to be an ID of a security group in the same tenant. Changing this creates a new security group rule. **Note**: Only one of `remote_group_id` or `remote_ip_prefix` may be set.

- `remote_ip_prefix` **String** (*Optional*) The remote CIDR, the value needs to be a valid CIDR (i.e. 192.168.0.0/16). Changing this creates a new security group rule. **Note**: Only one of `remote_group_id` or `remote_ip_prefix` may be set.

- `sdn` **String** (*Optional*) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".


## Attributes Reference
- `direction` **String** See Argument Reference above.

- `ethertype` **String** See Argument Reference above.

- `security_group_id` **String** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `port_range_max` **Number** See Argument Reference above.

- `port_range_min` **Number** See Argument Reference above.

- `protocol` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `remote_group_id` **String** See Argument Reference above.

- `remote_ip_prefix` **String** See Argument Reference above.

- `sdn` **String** See Argument Reference above.

- `id` **String** ID of the resource.



## Import

Security Group Rules can be imported using the `id`, e.g.

```shell
terraform import vkcs_networking_secgroup_rule.secgroup_rule_1 aeb68ee3-6e9d-4256-955c-9584a6212745
```
