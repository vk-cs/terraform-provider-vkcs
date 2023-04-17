---
layout: "vkcs"
page_title: "vkcs: vkcs_networking_secgroup"
description: |-
  Manages a security group resource within VKCS.
---

# vkcs_networking_secgroup

Manages a security group resource within VKCS.

## Example Usage
```terraform
resource "vkcs_networking_secgroup" "secgroup_1" {
  name        = "secgroup_1"
  description = "My security group"
}
```

## Argument Reference
- `name` **String** (***Required***) A unique name for the security group.

- `delete_default_rules` **Boolean** (*Optional*) Whether or not to delete the default egress security rules. This is `false` by default. See the below note for more information.

- `description` **String** (*Optional*) A unique name for the security group.

- `region` **String** (*Optional*) The region in which to obtain the networking client. A networking client is needed to create a port. If omitted, the `region` argument of the provider is used. Changing this creates a new security group.

- `sdn` **String** (*Optional*) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `tags` <strong>Set of </strong>**String** (*Optional*) A set of string tags for the security group.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `all_tags` <strong>Set of </strong>**String** The collection of tags assigned on the security group, which have been explicitly and implicitly added.

- `id` **String** ID of the resource.



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

```shell
terraform import vkcs_networking_secgroup.secgroup_1 38809219-5e8a-4852-9139-6f461c90e8bc
```
