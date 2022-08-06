---
layout: "vkcs"
page_title: "vkcs: kubernetes_node_group"
description: |-
  Get information on clusters node group.
---

# vkcs\_kubernetes\_node\_group

Provides a cluster node group resource. This can be used to create, modify and delete cluster's node group.

## Example Usage
```
resource "vkcs_kubernetes_node_group" "mynodegroup" {
    cluster_id = your_cluster_id
    name = my_new_node_group
    availability_zones = ["MS1"]
}
```

## Argument Reference

The following arguments are supported:

* `autoscaling_enabled` - (Optional) Determines whether the autoscaling is enabled.
* `availability_zones` - (Optional) The list of availability zones of the node group.
  Zones `MS1` and  `GZ1` are available. By default, node group is being created at
  cluster's zone.
  **Important:** Receiving default AZ add it manually to your main.tf config to sync it with state
  to avoid node groups force recreation in the future.
* `cluster_id` - (Required) The UUID of the existing cluster.
* `flavor_id` - (Optional) The flavor UUID of this node group.
* `labels` - (Optional) The list of objects representing representing additional
  properties of the node group. Each object should have attribute "key".
  Object may also have optional attribute "value".
* `max_nodes` - (Optional) The maximum allowed nodes for this node group.
* `min_nodes` - (Optional) The minimum allowed nodes for this node group. Default to 0 if not set.
* `name` - (Required) The name of node group to create.
 Changing this will force to create a new node group.
* `node_count` - (Required) The node count for this node group. Should be greater than 0.
 If `autoscaling_enabled` parameter is set, this attribute will be ignored during update.
* `taints` - (Optional) The list of objects representing node group taints. Each
  object should have following attributes: key, value, effect.
* `volume_size` - (Optional) The size in GB for volume to load nodes from.
 Changing this will force to create a new node group.
* `volume_type` - (Optional) The volume type to load nodes from.
 Changing this will force to create a new node group.
* `max_node_unavailable` - (Optional) The maximum number of nodes that can fail during an upgrade. The default value is 25 percent.


## Attributes
`id` is set to the ID of the found cluster template. In addition, the following
attributes are exported:

* `autoscaling_enabled` - Determines whether the autoscaling is enabled.
* `availability_zones` - The list of availability zones of the node group.
* `cluster_id` - The UUID of cluster that node group belongs.
* `flavor_id` - The UUID of a flavor.
* `labels` - The list of key value pairs representing additional
  properties of the node group.
* `max_nodes` - The maximum amount of nodes in node group.
* `min_nodes` - The minimum amount of nodes in node group.
* `name` - The name of the node group.
* `node_count` - The count of nodes in node group.
* `nodes` - The list of node group's node objects.
* `state` - Determines current state of node group (RUNNING, SHUTOFF, ERROR).
* `taints` - The list of objects representing node group taints.
* `uuid` - The UUID of the cluster's node group.
* `volume_size` - The size in GB for volume to load nodes from.
* `volume_type` - The volume type to load nodes from.
* `max_node_unavailable` - The maximum number of nodes that can fail during an upgrade.

## Import

Node groups can be imported using the `id`, e.g.

```
$ terraform import vkcs_kubernetes_node_group.ng ng_uuid
```
