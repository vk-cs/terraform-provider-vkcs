---
subcategory: "Firewall"
layout: "vkcs"
page_title: "vkcs: vkcs_networking_secgroup_rule"
description: |-
  Manages a security group rule resource within VKCS.
---

# vkcs_networking_secgroup_rule

Manages a security group rule resource within VKCS.

## Example Usage
### Basic-rule one protocol one port
```terraform
resource "vkcs_networking_secgroup_rule" "etcd_app_clients" {
  description       = "etcd app clients rule"
  security_group_id = vkcs_networking_secgroup.etcd.id
  direction         = "ingress"
  protocol          = "tcp"
  port_range_max    = 2379
  port_range_min    = 2379
  remote_ip_prefix  = vkcs_networking_subnet.app.cidr
}
```

### Rule for all ports for particular protocol (udp)
```terraform
resource "vkcs_networking_secgroup_rule" "all_udp" {
  description       = "All inbound UDP traffic from etcd hosts"
  security_group_id = vkcs_networking_secgroup.etcd.id
  direction         = "ingress"
  protocol          = "udp"
  remote_group_id   = vkcs_networking_secgroup.etcd.id
}
```

### Rule for all protocols
```terraform
resource "vkcs_networking_secgroup_rule" "all" {
  description       = "Any inbound traffic from etcd hosts"
  security_group_id = vkcs_networking_secgroup.etcd.id
  direction         = "ingress"
  remote_group_id   = vkcs_networking_secgroup.etcd.id
}
```

## Argument Reference
- `direction` **required** *string* &rarr;  The direction of the rule, valid values are __ingress__ or __egress__. Changing this creates a new security group rule.

- `security_group_id` **required** *string* &rarr;  The security group id the rule should belong to, the value needs to be an ID of a security group in the same tenant. Changing this creates a new security group rule.

- `description` optional *string* &rarr;  A description of the rule. Changing this creates a new security group rule.

- `ethertype` optional deprecated *string* &rarr;  The layer 3 protocol type, the only valid value is __IPv4__. Changing this creates a new security group rule. **Deprecated** Only IPv4 can be used as ethertype. This argument is deprecated, please do not use it.

- `port_range_max` optional *number* &rarr;  The higher part of the allowed port range, valid integer value needs to be between 1 and 65535. To specify all ports, `port_range_min` and `port_range_max` arguments must be absent. Changing this creates a new security group rule.

- `port_range_min` optional *number* &rarr;  The lower part of the allowed port range, valid integer value needs to be between 1 and 65535. To specify all ports, `port_range_min` and `port_range_max` arguments must be absent. Changing this creates a new security group rule.

- `protocol` optional *string* &rarr;  The layer 4 protocol type, valid values are following. Changing this creates a new security group rule. This is required if you want to specify a port range.
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

- `region` optional *string* &rarr;  The region in which to obtain the networking client. A networking client is needed to create a port. If omitted, the `region` argument of the provider is used. Changing this creates a new security group rule.

- `remote_group_id` optional *string* &rarr;  The remote group id, the value needs to be an ID of a security group in the same tenant. Changing this creates a new security group rule. <br>**Note:** Only one of `remote_group_id` or `remote_ip_prefix` may be set.

- `remote_ip_prefix` optional *string* &rarr;  The remote CIDR, the value needs to be a valid CIDR (i.e. 192.168.0.0/16). Changing this creates a new security group rule. <br>**Note:** Only one of `remote_group_id` or `remote_ip_prefix` may be set.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is project's default SDN.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Security Group Rules can be imported using the `id`, e.g.

```shell
terraform import vkcs_networking_secgroup_rule.secgroup_rule_1 aeb68ee3-6e9d-4256-955c-9584a6212745
```
