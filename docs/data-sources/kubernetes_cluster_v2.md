---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_cluster_v2"
description: |-
  Retrieves information about an existing next-generation Kubernetes cluster.
---

# vkcs_kubernetes_cluster_v2



## Example Usage
```terraform
data "vkcs_kubernetes_cluster_v2" "k8s_cluster" {
  id = vkcs_kubernetes_cluster_v2.k8s_cluster.id
}
```

## Argument Reference
- `id` **required** *string* &rarr;  The cluster's ID. Used to fetch a specific cluster.

- `region` optional *string* &rarr;  The region in which the cluster resides. If omitted, the provider's configured `region` is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `api_address` *string* &rarr;  The full URL of the Kubernetes API server (e.g., `https://api.example.com:6443`).

- `api_lb_fip` *string* &rarr;  The floating IP address of the API load balancer (present only when `public_ip` is `true`).

- `api_lb_vip` *string* &rarr;  The internal virtual IP address of the API load balancer.

- `availability_zones` *set of* *string* &rarr;  The availability zones in which the master nodes reside.

- `cluster_type` *string* &rarr;  The cluster type: `standard` (single AZ) or `regional` (three AZs).

- `created_at` *string* &rarr;  The timestamp when the cluster was created (ISO 8601 format).

- `description` *string* &rarr;  A user-provided description of the cluster.

- `external_network_id` *string* &rarr;  The ID of the external network providing internet access to the cluster.

- `insecure_registries` *set of* *string* &rarr;  A list of container registry addresses that can be accessed without certificate validation.

- `k8s_config` *string* &rarr;  The raw kubeconfig file contents used to authenticate to the Kubernetes API.

- `labels` *map of* *string* &rarr;  Kubernetes labels applied to the entire cluster.

- `loadbalancer_allowed_cidrs` *set of* *string* &rarr;  CIDR blocks permitted to access the cluster's API load balancer.

- `loadbalancer_subnet_id` *string* &rarr;  The ID of the subnet where LoadBalancer type services are provisioned.

- `master_count` *number* &rarr;  The number of master nodes (1, 3, or 5) ensuring etcd quorum.

- `master_disks`  *set* &rarr;  Information about the disks attached to master nodes (root, etcd, etcd-events). Each object contains `type` (disk type) and `size` (GB).
    - `size` *number* &rarr;  The allocated size of the disk in gigabytes (GB).

    - `type` *string* &rarr;  The storage volume type used for this disk (e.g., `ceph-ssd`, `ceph-hdd`). Corresponds to the available volume types in the region.


- `master_flavor` *string* &rarr;  The flavor ID (VM type) used for master nodes.

- `name` *string* &rarr;  The user-assigned name of the cluster.

- `network_id` *string* &rarr;  The ID of the VPC network where the cluster is deployed.

- `network_plugin` *string* &rarr;  The CNI network plugin in use (always `calico`).

- `node_groups`  *set* &rarr;  A list of node groups associated with the cluster. Each object includes `id`, `name`, `flavor`, `node_count`, and `availability_zone`.
    - `availability_zone` *string* &rarr;  The availability zone where the nodes of this group are deployed.

    - `flavor` *string* &rarr;  The flavor ID (VM type) used for each worker node in this group.

    - `id` *string* &rarr;  The unique identifier of the node group.

    - `name` *string* &rarr;  The user-assigned name of the node group.

    - `node_count` *number* &rarr;  The current number of worker nodes in the group. This value reflects the actual node count, accounting for both manual scaling and autoscaling actions.


- `pods_ipv4_cidr` *string* &rarr;  The IPv4 CIDR block allocated for pod networking.

- `project_id` *string* &rarr;  The VK Cloud project ID owning the cluster.

- `public_ip` *boolean* &rarr;  Indicates whether a floating IP is assigned to the cluster's API endpoint.

- `status` *string* &rarr;  The current operational status of the cluster (e.g., `RUNNING`, `RECONCILING`).

- `subnet_id` *string* &rarr;  The ID of the subnet used by the cluster's nodes.

- `uuid` *string* &rarr;  The cluster's UUID.

- `version` *string* &rarr;  The installed Kubernetes version (e.g., `v1.34.2`).


