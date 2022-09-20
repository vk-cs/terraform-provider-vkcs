---
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_node_group"
description: |-
  Get information on clusters node group.
---

# vkcs_kubernetes_node_group

Use this data source to get the ID of an available VKCS kubernetes clusters node group.

## Example Usage
```terraform
data "vkcs_kubernetes_node_group" "mynodegroup" {
  uuid = "mynguuid"
}
```
## Argument Reference
- `uuid` **String** (***Required***) The UUID of the cluster's node group.

- `autoscaling_enabled` **Boolean** (*Optional*) Determines whether the autoscaling is enabled.

- `flavor_id` **String** (*Optional*) The id of flavor.

- `max_node_unavailable` **Number** (*Optional*) Specified as a percentage. The maximum number of nodes that can fail during an upgrade.

- `max_nodes` **Number** (*Optional*) The maximum amount of nodes in node group.

- `min_nodes` **Number** (*Optional*) The minimum amount of nodes in node group.

- `name` **String** (*Optional*) The name of the node group.

- `node_count` **Number** (*Optional*) The count of nodes in node group.

- `volume_size` **Number** (*Optional*) The amount of memory of volume in Gb.

- `volume_type` **String** (*Optional*) The type of volume.


## Attributes Reference
- `uuid` **String** See Argument Reference above.

- `autoscaling_enabled` **Boolean** See Argument Reference above.

- `flavor_id` **String** See Argument Reference above.

- `max_node_unavailable` **Number** See Argument Reference above.

- `max_nodes` **Number** See Argument Reference above.

- `min_nodes` **Number** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `node_count` **Number** See Argument Reference above.

- `volume_size` **Number** See Argument Reference above.

- `volume_type` **String** See Argument Reference above.

- `availability_zones` **String** The list of availability zones of the node group.

- `cluster_id` **String** The UUID of cluster that node group belongs.

- `id` **String** ID of the resource.

- `nodes` **Object** The list of node group's node objects.

- `state` **String** Determines current state of node group (RUNNING, SHUTOFF, ERROR).


