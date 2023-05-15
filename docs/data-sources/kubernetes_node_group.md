---
subcategory: "Kubernetes"
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
- `uuid` **required** *string* &rarr;  The UUID of the cluster's node group.

- `autoscaling_enabled` optional *boolean* &rarr;  Determines whether the autoscaling is enabled.

- `flavor_id` optional *string* &rarr;  The id of flavor.

- `max_node_unavailable` optional *number* &rarr;  Specified as a percentage. The maximum number of nodes that can fail during an upgrade.

- `max_nodes` optional *number* &rarr;  The maximum amount of nodes in node group.

- `min_nodes` optional *number* &rarr;  The minimum amount of nodes in node group.

- `name` optional *string* &rarr;  The name of the node group.

- `node_count` optional *number* &rarr;  The count of nodes in node group.

- `volume_size` optional *number* &rarr;  The amount of memory of volume in Gb.

- `volume_type` optional *string* &rarr;  The type of volume.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `availability_zones` *string* &rarr;  The list of availability zones of the node group.

- `cluster_id` *string* &rarr;  The UUID of cluster that node group belongs.

- `id` *string* &rarr;  ID of the resource.

- `nodes` *object* &rarr;  The list of node group's node objects.

- `state` *string* &rarr;  Determines current state of node group (RUNNING, SHUTOFF, ERROR).


