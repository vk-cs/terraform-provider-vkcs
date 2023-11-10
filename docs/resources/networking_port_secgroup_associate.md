---
subcategory: "Network"
layout: "vkcs"
page_title: "vkcs: vkcs_networking_port_secgroup_associate"
description: |-
  Manages a port's security groups within VKCS.
---

# vkcs_networking_port_secgroup_associate

Manages a port's security groups within VKCS. Useful, when the port was created not by Terraform. It should not be used, when the port was created directly within Terraform.

When the resource is deleted, Terraform doesn't delete the port, but unsets the list of user defined security group IDs.  However, if `enforce` is set to `true` and the resource is deleted, Terraform will remove all assigned security group IDs.

## Example Usage
### Append a security group to an existing port
```terraform
resource "vkcs_networking_port_secgroup_associate" "append_secgroups" {
  port_id = vkcs_networking_port.persistent_etcd.id
  security_group_ids = [
    vkcs_networking_secgroup.etcd.id,
  ]
}
```

### Enforce a security group to an existing port
```terraform
resource "vkcs_networking_port_secgroup_associate" "replace_secgroups" {
  port_id = vkcs_networking_port.persistent_etcd.id
  security_group_ids = [
    vkcs_networking_secgroup.http.id,
  ]
  enforce = true
}
```

### Remove all security groups from an existing port
```terraform
resource "vkcs_networking_port_secgroup_associate" "remove_secgroups" {
  port_id = vkcs_networking_port.persistent_etcd.id
  security_group_ids = []
  enforce = true
}
```
## Argument Reference
- `port_id` **required** *string* &rarr;  An UUID of the port to apply security groups to.

- `security_group_ids` **required** *set of* *string* &rarr;  A list of security group IDs to apply to the port. The security groups must be specified by ID and not name (as opposed to how they are configured with the Compute Instance).

- `enforce` optional *boolean* &rarr;  Whether to replace or append the list of security groups, specified in the `security_group_ids`. Defaults to `false`.

- `region` optional *string* &rarr;  The region in which to obtain the networking client. A networking client is needed to manage a port. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is project's default SDN.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `all_security_group_ids` *set of* *string* &rarr;  The collection of Security Group IDs on the port which have been explicitly and implicitly added.

- `id` *string* &rarr;  ID of the resource.


