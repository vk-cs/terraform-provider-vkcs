---
subcategory: "Firewall"
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
- `description` optional *string* &rarr;  Human-readable description the the subnet.

- `name` optional *string* &rarr;  The name of the security group.

- `region` optional *string* &rarr;  The region in which to obtain the Network client. A Network client is needed to retrieve security groups ids. If omitted, the `region` argument of the provider is used.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `secgroup_id` optional *string* &rarr;  The ID of the security group.

- `tags` optional *set of* *string* &rarr;  The list of security group tags to filter.

- `tenant_id` optional *string* &rarr;  The owner of the security group.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `all_tags` *set of* *string* &rarr;  The set of string tags applied on the security group.

- `id` *string* &rarr;  ID of the found security group.


