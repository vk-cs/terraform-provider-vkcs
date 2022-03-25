---
layout: "vkcs"
page_title: "vkcs: networking_secgroup"
description: |-
  Get information on an VKCS Security Group.
---

# vkcs\_networking\_secgroup

Use this data source to get the ID of an available VKCS security group.

## Example Usage

```hcl
data "vkcs_networking_secgroup" "secgroup" {
  name = "tf_test_secgroup"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the Network client.
  A Network client is needed to retrieve security groups ids. If omitted, the
  `region` argument of the provider is used.

* `secgroup_id` - (Optional) The ID of the security group.

* `name` - (Optional) The name of the security group.

* `description` - (Optional) Human-readable description the the subnet.

* `tags` - (Optional) The list of security group tags to filter.

* `tenant_id` - (Optional) The owner of the security group.

* `sdn` - (Optional) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

## Attributes Reference

`id` is set to the ID of the found security group. In addition, the following
attributes are exported:

* `name` - See Argument Reference above.
* `description`- See Argument Reference above.
* `all_tags` - The set of string tags applied on the security group.
* `region` - See Argument Reference above.
* `sdn` - See Argument Reference above.
