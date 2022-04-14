---
layout: "vkcs"
page_title: "vkcs: compute_servergroup"
description: |-
  Manages a Server Group resource within VKCS.
---

# vkcs\_compute\_servergroup

Manages a Server Group resource within VKCS.

## Example Usage

```hcl
resource "vkcs_compute_servergroup" "test-sg" {
  name     = "my-sg"
  policies = ["anti-affinity"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the Compute client.
    If omitted, the `region` argument of the provider is used. Changing
    this creates a new server group.

* `name` - (Required) A unique name for the server group. Changing this creates
    a new server group.

* `policies` - (Required) The set of policies for the server group. All policies
    are mutually exclusive. See the Policies section for more information.
    Changing this creates a new server group.

* `value_specs` - (Optional) Map of additional options.

## Policies

* `affinity` - All instances/servers launched in this group will be hosted on
    the same compute node.

* `anti-affinity` - All instances/servers launched in this group will be
    hosted on different compute nodes.

* `soft-affinity` - All instances/servers launched in this group will be hosted
    on the same compute node if possible, but if not possible they
    still will be scheduled instead of failure.

* `soft-anti-affinity` - All instances/servers launched in this group will be
    hosted on different compute nodes if possible, but if not possible they
    still will be scheduled instead of failure.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `policies` - See Argument Reference above.
* `members` - The instances that are part of this server group.

## Import

Server Groups can be imported using the `id`, e.g.

```
$ terraform import vkcs_compute_servergroup.test-sg 1bc30ee9-9d5b-4c30-bdd5-7f1e663f5edf
```
