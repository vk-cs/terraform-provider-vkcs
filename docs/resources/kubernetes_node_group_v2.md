---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_node_group_v2"
description: |-
  Provides a Kubernetes cluster node group resource. This can be used to create, modify, and delete Kubernetes cluster node group.
---

# vkcs_kubernetes_node_group_v2

Provides a Kubernetes cluster node group resource. This can be used to create, modify, and delete Kubernetes cluster node group.

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

- `disk_size` **required** *number* &rarr;  The size of the root volume, in gigabytes (GB). Must be at least 1. **Forces replacement** on change.

- `disk_type` **required** *string* &rarr;  The root volume type. For example: `ceph-ssd`. Use the `vkcs_kubernetes_volume_types_v2` data source to list available types for the selected availability zone. **Forces replacement** on change.

- `name` **required** *string* &rarr;  A unique name for the node group within the cluster. Must be 3-25 characters, lowercase alphanumeric and hyphens only. **Forces replacement** on change.

- `node_flavor` **required** *string* &rarr;  The flavor ID used for worker nodes. Changing this triggers a rolling upgrade of the nodes.

- `parallel_upgrade_chunk` **required** *number* &rarr;  The maximum percentage of nodes (1-100) that can be upgraded simultaneously during a rolling upgrade.

- `scale_type` **required** *string* &rarr;  Type of scaling for the node group. Must be either `fixed_scale` or `auto_scale`. If `scale_type` is `auto_scale`, the condition `auto_scale_min_size` <= `auto_scale_max_size` must be met.

- `auto_scale_max_size` optional *number* &rarr;  The maximum number of nodes the autoscaler may scale out to. Must be greater than or equal to 0 and also greater than or equal to `auto_scale_min_size`. Required if scale_type is `auto_scale`.

- `auto_scale_min_size` optional *number* &rarr;  The minimum number of nodes the autoscaler can scale down to. Must be at least 0. Required when `scale_type` is `auto_scale`.

- `fixed_scale_node_count` optional *number* &rarr;  The desired number of nodes. Minimum value is 0. This argument is required when `scale_type` is `fixed_scale`.

- `labels` optional *map of* *string* &rarr;  Kubernetes labels to apply to the nodes in this node group. Both keys and values must conform to Kubernetes label syntax.

- `region` optional *string* &rarr;  The region where the node group will be created. If omitted, the provider's `region` is used. **Forces replacement** on change.

- `taints`  *set* &rarr;  Kubernetes taints to apply to the nodes in this node group. Each taint must specify `key`, `value`, and `effect` — where `effect` is one of: `NoSchedule`, `PreferNoSchedule`, `NoExecute`.
    - `effect` **required** *string* &rarr;  The effect of the taint. Must be one of: `NoSchedule`, `PreferNoSchedule`, `NoExecute`.

    - `key` **required** *string* &rarr;  The key of the taint. Must conform to Kubernetes label key syntax.

    - `value` **required** *string* &rarr;  The value of the taint. Must conform to Kubernetes label value syntax.



## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `auto_scale_node_count` *number* &rarr;  The current number of nodes in the node group. Only applicable when `scale_type` is `auto_scale`. This value is computed by the provider and cannot be set by the user.

- `created_at` *string* &rarr;  The timestamp when the node group was created, in ISO 8601 format.

- `id` *string* &rarr;  The unique identifier of the node group.

- `uuid` *string* &rarr;  The UUID of the node group. It is generated automatically.



## Import

Node groups can be imported using the `id`, e.g.

```shell
terraform import vkcs_kubernetes_node_group_v2.k8s_node_group 39WHkEHwtXy1YWqka4D5xuBJxw4
```
