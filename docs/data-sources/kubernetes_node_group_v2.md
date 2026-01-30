---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_node_group_v2"
description: |-
  Retrieves information about a specific node group within a next-generation Kubernetes cluster.
---

# vkcs_kubernetes_node_group_v2



## Example Usage
```terraform
data "vkcs_kubernetes_node_group_v2" "node_group" {
  id = vkcs_kubernetes_node_group_v2.k8s_node_group.id
}
```

## Argument Reference
- `id` **required** *string* &rarr;  The ID of the node group to retrieve.

- `region` optional *string* &rarr;  The region where the node group resides. Defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `auto_scale_max_size` *number* &rarr;  The maximum allowed node count for `auto_scale`.

- `auto_scale_min_size` *number* &rarr;  The minimum allowed node count for `auto_scale`.

- `auto_scale_node_count` *number* &rarr;  The current number of nodes when `scale_type` is `auto_scale`.

- `availability_zone` *string* &rarr;  The availability zone where the nodes are located.

- `cluster_id` *string* &rarr;  The ID of the parent Kubernetes cluster.

- `created_at` *string* &rarr;  The creation timestamp of the node group (ISO 8601 format).

- `disk_size` *number* &rarr;  The size of the root volume in GB.

- `disk_type` *string* &rarr;  The type of root volume attached to each node (e.g., `ceph-ssd`).

- `fixed_scale_node_count` *number* &rarr;  The desired number of nodes when `scale_type` is `fixed_scale`.

- `labels` *map of* *string* &rarr;  Kubernetes labels applied to the nodes in this group.

- `name` *string* &rarr;  The user-assigned name of the node group.

- `node_flavor` *string* &rarr;  The flavor ID (VM type) of each node in the group.

- `parallel_upgrade_chunk` *number* &rarr;  The maximum percentage of nodes that can be upgraded simultaneously (1-100).

- `scale_type` *string* &rarr;  The scaling strategy: `fixed_scale` (static node count) or `auto_scale` (dynamic based on load).

- `taints`  *set* &rarr;  Taints applied to the nodes. Each taint contains `key`, `value`, and `effect` (e.g., `NoSchedule`).
    - `effect` *string* &rarr;  The effect of the taint. Allowed values are: `NoSchedule`, `PreferNoSchedule`, `NoExecute`.

    - `key` *string* &rarr;  The taint key. Must conform to Kubernetes label key syntax (e.g., `node.kubernetes.io/unreachable`).

    - `value` *string* &rarr;  The taint value. Must conform to Kubernetes label value syntax (e.g., `production`).


- `uuid` *string* &rarr;  The node group's UUID


