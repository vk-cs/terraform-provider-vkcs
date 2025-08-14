---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_node_group"
description: |-
  Get information on clusters node group.
---

# vkcs_kubernetes_node_group

Use this data source to get information on VKCS Kubernetes cluster's node group.

## Example Usage
```terraform
data "vkcs_kubernetes_node_group" "k8s_node_group" {
  id = vkcs_kubernetes_node_group.default_ng.id
}
```
## Argument Reference
- `autoscaling_enabled` optional *boolean* &rarr;  Determines whether the autoscaling is enabled.

- `flavor_id` optional *string* &rarr;  The id of the flavor.

- `id` optional *string* &rarr;  The UUID of the cluster's node group.

- `max_node_unavailable` optional *number* &rarr;  Specified as a percentage. The maximum number of nodes that can fail during an upgrade.

- `max_nodes` optional *number* &rarr;  The maximum amount of nodes in the node group.

- `min_nodes` optional *number* &rarr;  The minimum amount of nodes in the node group.

- `name` optional *string* &rarr;  The name of the node group.

- `node_count` optional *number* &rarr;  The count of nodes in the node group.

- `region` optional *string* &rarr;  The region to obtain the service client. If omitted, the `region` argument of the provider is used.<br>**New since v0.4.0**.

- `uuid` optional deprecated *string* &rarr;  The UUID of the cluster's node group. **Deprecated** This argument is deprecated, please, use the `id` attribute instead.

- `volume_size` optional *number* &rarr;  The amount of memory in the volume in GB

- `volume_type` optional *string* &rarr;  The type of the volume.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `availability_zones` *string* &rarr;  The list of availability zones of the node group.

- `cluster_id` *string* &rarr;  The UUID of cluster that node group belongs.

- `nodes`  *list* &rarr;  The list of node group's node objects.
    - `created_at` *string* &rarr;  Time when a node was created.

    - `name` *string* &rarr;  Name of a node.

    - `node_group_id` *string* &rarr;  The node group id.

    - `updated_at` *string* &rarr;  Time when a node was updated.

    - `uuid` *string* &rarr;  UUID of a node.


- `state` *string* &rarr;  Determines current state of node group (RUNNING, SHUTOFF, ERROR).


