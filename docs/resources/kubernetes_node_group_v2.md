---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_node_group_v2"
description: |-
  Creates and manages a node group within an existing next-generation Kubernetes cluster.
---

# vkcs_kubernetes_node_group_v2



```terraform
resource "vkcs_kubernetes_node_group_v2" "k8s_node_group" {
  cluster_id = vkcs_kubernetes_cluster_v2.k8s_cluster.id
  name       = "k8s-node-group"

  node_flavor       = data.vkcs_compute_flavor.basic.id
  availability_zone = "MS1"

  scale_type             = "fixed_scale"
  fixed_scale_node_count = 3

  parallel_upgrade_chunk = 40

  disk_type = "high-iops"
  disk_size = 20

  depends_on = [
    vkcs_kubernetes_cluster_v2.k8s_cluster,
  ]
}
```

## Argument Reference
- `availability_zone` **required** *string* &rarr;  The availability zone where the nodes will be placed. **Forces replacement** on change.

- `cluster_id` **required** *string* &rarr;  The ID of the parent cluster. **Forces replacement** on change.

- `disk_size` **required** *number* &rarr;  The size of the root volume in GB. Minimum is 1 GB. **Forces replacement** on change.

- `disk_type` **required** *string* &rarr;  The type of root volume (e.g., `ceph-ssd`). Use `vkcs_kubernetes_volume_types_v2` to list available types. **Forces replacement** on change.

- `name` **required** *string* &rarr;  A unique name for the node group within the cluster. Must be 3-25 characters, lowercase alphanumeric and hyphens only. **Forces replacement** on change.

- `node_flavor` **required** *string* &rarr;  The flavor ID for each node. Changing this triggers a rolling upgrade of the nodes.

- `parallel_upgrade_chunk` **required** *number* &rarr;  The maximum percentage of nodes (1-100) that can be upgraded simultaneously during a rolling upgrade.

- `scale_type` **required** *string* &rarr;  Type of scaling for the node group. Must be either `fixed_scale` or `auto_scale`. If `scale_type` is `auto_scale`, the condition `auto_scale_min_size` <= `auto_scale_max_size` must be met.

- `auto_scale_max_size` optional *number* &rarr;  The maximum number of nodes the autoscaler may scale out to. Must be greater than or equal to 0 and also greater than or equal to `auto_scale_min_size`. Required if scale_type is `auto_scale`.

- `auto_scale_min_size` optional *number* &rarr;  The minimum number of nodes the autoscaler may scale down to. Must be greater than or equal to 0. Required if `scale_type` is `auto_scale`.

- `fixed_scale_node_count` optional *number* &rarr;  The desired number of nodes. Minimum value is 0. Required if `scale_type` is `fixed_scale`.

- `labels` optional *map of* *string* &rarr;  Kubernetes labels to apply to nodes in this group. Keys/values must be valid Kubernetes label strings.

- `region` optional *string* &rarr;  The region where the node group will be created. Defaults to provider's `region`. **Forces replacement** on change.

- `taints`  *set* &rarr;  Taints to apply to nodes. Each taint must have `key`, `value`, and `effect` (one of `NoSchedule`, `PreferNoSchedule`, `NoExecute`).
    - `effect` **required** *string* &rarr;  The effect of the taint. Must be one of: `NoSchedule`, `PreferNoSchedule`, `NoExecute`.

    - `key` **required** *string* &rarr;  The taint key. Must be a valid Kubernetes label key.

    - `value` **required** *string* &rarr;  The taint value. Must be a valid Kubernetes label value.



## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `auto_scale_node_count` *number* &rarr;  During the cluster lifecycle, it indicates the current number of nodes in the node group if `scale_type` is `auto_scale`.

- `created_at` *string* &rarr;  The timestamp when the node group was created (ISO 8601 format).

- `id` *string* &rarr;  The unique identifier of the node group.

- `uuid` *string* &rarr;  The node group's UUID. It is generated automatically.



## Import

Node groups can be imported using the `id`, e.g.

```shell
terraform import vkcs_kubernetes_node_group_v2.k8s_node_group 39WHkEHwtXy1YWqka4D5xuBJxw4
```
