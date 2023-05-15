---
subcategory: "Virtual Machines"
layout: "vkcs"
page_title: "vkcs: vkcs_compute_servergroup"
description: |-
  Manages a Server Group resource within VKCS.
---

# vkcs_compute_servergroup

Manages a Server Group resource within VKCS.

## Example Usage
```terraform
resource "vkcs_compute_servergroup" "test-sg" {
  name     = "my-sg"
  policies = ["anti-affinity"]
}
```
## Argument Reference
- `name` **required** *string* &rarr;  A unique name for the server group. Changing this creates a new server group.

- `policies` optional *string* &rarr;  The set of policies for the server group. All policies are mutually exclusive. See the Policies section for more information. Changing this creates a new server group.

- `region` optional *string* &rarr;  The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used. Changing this creates a new server group.

- `value_specs` optional *map of* *string* &rarr;  Map of additional options.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `members` *string* &rarr;  The instances that are part of this server group.


## Policies

* `affinity` - All instances/servers launched in this group will be hosted on the same compute node.

* `anti-affinity` - All instances/servers launched in this group will be hosted on different compute nodes.

* `soft-affinity` - All instances/servers launched in this group will be hosted on the same compute node if possible, but if not possible they still will be scheduled instead of failure.

* `soft-anti-affinity` - All instances/servers launched in this group will be hosted on different compute nodes if possible, but if not possible they still will be scheduled instead of failure.

## Import

Server Groups can be imported using the `id`, e.g.
```shell
terraform import vkcs_compute_servergroup.test-sg 1bc30ee9-9d5b-4c30-bdd5-7f1e663f5edf
```
