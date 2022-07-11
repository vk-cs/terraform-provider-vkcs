---
layout: "vkcs"
page_title: "vkcs: kubernetes_node_group"
description: |-
  Get information on clusters node group.
---

# VKCS Kubernetes Node Group

Use this data source to get the ID of an available VKCS kubernetes clusters node group.

## Example Usage
```hcl
data "vkcs_kubernetes_node_group" "mynodegroup" {
  uuid = "mynguuid"
}
```

## Argument Reference

The following arguments are supported:

* `uuid` - (Required) The UUID of the cluster's node group.

    
## Attributes
`id` is set to the ID of the found cluster template. In addition, the following
attributes are exported:

* `autoscaling_enabled` - Determines whether the autoscaling is enabled.
* `availability_zones` - The list of availability zones of the node group.
* `cluster_id` - The UUID of cluster that node group belongs.
* `flavor_id` - The id of flavor.
* `max_nodes` - The maximum amount of nodes in node group.
* `min_nodes` - The minimum amount of nodes in node group.
* `name` - The name of the node group.
* `node_count` - The count of nodes in node group.
* `nodes` - The list of node group's node objects.
* `state` - Determines current state of node group (RUNNING, SHUTOFF, ERROR).
* `uuid` - The UUID of the cluster's node group.
* `volume_size` - The amount of memory of volume in Gb.
* `volume_type` - The type of volume.
* `max_node_unavailable` - Specified as a percentage. The maximum number of nodes that can fail during an upgrade.
