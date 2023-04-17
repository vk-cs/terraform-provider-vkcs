---
layout: "vkcs"
page_title: "vkcs: vkcs_networking_secgroup"
description: |-
  Get information on an VKCS Security Group.
---

# vkcs_networking_secgroup

Use this data source to get the ID of an available VKCS security group.

## Example Usage

```terraform
data "vkcs_networking_secgroup" "secgroup" {
  name = "tf_test_secgroup"
}
```

## Argument Reference
- `description` **String** (*Optional*) Human-readable description the the subnet.

- `name` **String** (*Optional*) The name of the security group.

- `region` **String** (*Optional*) The region in which to obtain the Network client. A Network client is needed to retrieve security groups ids. If omitted, the `region` argument of the provider is used.

- `sdn` **String** (*Optional*) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `secgroup_id` **String** (*Optional*) The ID of the security group.

- `tags` <strong>Set of </strong>**String** (*Optional*) The list of security group tags to filter.

- `tenant_id` **String** (*Optional*) The owner of the security group.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `all_tags` <strong>Set of </strong>**String** The set of string tags applied on the security group.

- `id` **String** ID of the found security group.


