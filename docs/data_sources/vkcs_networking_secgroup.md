---
layout: "vkcs"
page_title: "VKCS: networking_secgroup"
description: |-
  Get information on an OpenStack Security Group.
---

# vkcs\_networking\_secgroup

Use this data source to get the ID of an available OpenStack security group.

## Example Usage

```hcl
data "vkcs_networking_secgroup" "secgroup" {
  name = "tf_test_secgroup"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve security groups ids. If omitted, the
  `region` argument of the provider is used.

* `secgroup_id` - (Optional) The ID of the security group.

* `name` - (Optional) The name of the security group.

* `description` - (Optional) Human-readable description the the subnet.

* `tags` - (Optional) The list of security group tags to filter.

* `tenant_id` - (Optional) The owner of the security group.

## Attributes Reference

`id` is set to the ID of the found security group. In addition, the following
attributes are exported:

* `name` - See Argument Reference above.
* `description`- See Argument Reference above.
* `all_tags` - The set of string tags applied on the security group.
* `region` - See Argument Reference above.
