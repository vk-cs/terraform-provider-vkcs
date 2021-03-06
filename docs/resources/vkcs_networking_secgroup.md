---
layout: "vkcs"
page_title: "vkcs: networking_secgroup"
description: |-
  Manages a security group resource within VKCS.
---

# vkcs\_networking\_secgroup

Manages a security group resource within VKCS.

## Example Usage

```hcl
resource "vkcs_networking_secgroup" "secgroup_1" {
  name        = "secgroup_1"
  description = "My security group"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the networking client.
    A networking client is needed to create a port. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    security group.

* `name` - (Required) A unique name for the security group.

* `description` - (Optional) A unique name for the security group.

* `delete_default_rules` - (Optional) Whether or not to delete the default
    egress security rules. This is `false` by default. See the below note
    for more information.

* `tags` - (Optional) A set of string tags for the security group.

* `sdn` - (Optional) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `all_tags` - The collection of tags assigned on the security group, which have
  been explicitly and implicitly added.
* `sdn` - See Argument Reference above.

## Default Security Group Rules

In most cases, VKCS will create some egress security group rules for each
new security group. These security group rules will not be managed by
Terraform, so if you prefer to have *all* aspects of your infrastructure
managed by Terraform, set `delete_default_rules` to `true` and then create
separate security group rules such as the following:

```hcl
resource "vkcs_networking_secgroup_rule" "secgroup_rule_v4" {
  direction         = "egress"
  ethertype         = "IPv4"
  security_group_id = "${vkcs_networking_secgroup.secgroup.id}"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_v6" {
  direction         = "egress"
  ethertype         = "IPv6"
  security_group_id = "${vkcs_networking_secgroup.secgroup.id}"
}
```

## Import

Security Groups can be imported using the `id`, e.g.

```
$ terraform import vkcs_networking_secgroup.secgroup_1 38809219-5e8a-4852-9139-6f461c90e8bc
```
