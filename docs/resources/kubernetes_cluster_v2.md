---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_cluster_v2"
description: |-
  Creates and manages a next-generation Kubernetes cluster.
---

# vkcs_kubernetes_cluster_v2



## Standard Kubernetes cluster
```terraform
resource "vkcs_kubernetes_cluster_v2" "k8s_cluster" {
  name        = "k8s-standard-cluster"
  description = "An example of a standard Kubernetes cluster v2 created via Terraform"
  version     = "v1.34.2"

  cluster_type       = "standard"
  master_count       = 1
  availability_zones = ["MS1"]
  master_flavor      = data.vkcs_compute_flavor.master.id

  network_id             = vkcs_networking_network.app.id
  subnet_id              = vkcs_networking_subnet.app.id
  loadbalancer_subnet_id = vkcs_networking_subnet.app.id
  network_plugin         = "calico"
  pods_ipv4_cidr         = "10.100.0.0/16"

  # If your configuration also defines a network for the instance,
  # ensure it is attached to a router before creating of the instance
  depends_on = [
    vkcs_networking_router_interface.app,
  ]
}
```

## Regional Kubernetes cluster
```terraform
resource "vkcs_kubernetes_cluster_v2" "k8s_cluster" {
  name        = "k8s-regional-cluster"
  description = "An example of a regional Kubernetes cluster v2 created via Terraform"
  version     = "v1.34.2"

  cluster_type       = "regional"
  master_count       = 3
  availability_zones = ["GZ1", "MS1", "ME1"]
  master_flavor      = data.vkcs_compute_flavor.master.id

  network_id             = vkcs_networking_network.app.id
  subnet_id              = vkcs_networking_subnet.app.id
  loadbalancer_subnet_id = vkcs_networking_subnet.app.id
  network_plugin         = "calico"
  pods_ipv4_cidr         = "10.100.0.0/16"

  # If your configuration also defines a network for the instance,
  # ensure it is attached to a router before creating of the instance
  depends_on = [
    vkcs_networking_router_interface.app,
  ]
}
```

## Argument Reference
- `availability_zones` **required** *set of* *string* &rarr;  A set of availability zones for the cluster's masters. For a `standard` cluster, provide 1, 3, or 5 zones. For a `regional` cluster, provide 3 or 5 zones. **Forces replacement** on change.

- `cluster_type` **required** *string* &rarr;  Either `standard` (single AZ) or `regional` (three AZs). **Forces replacement** on change.

- `loadbalancer_subnet_id` **required** *string* &rarr;  The ID of the subnet where LoadBalancer services will be allocated IPs. **Forces replacement** on change.

- `master_count` **required** *number* &rarr;  The number of master nodes. Must be an odd number to maintain etcd quorum (1, 3, or 5). **Forces replacement** on change.

- `master_flavor` **required** *string* &rarr;  The flavor ID for master nodes. Changing this triggers a rolling upgrade of the control plane.

- `name` **required** *string* &rarr;  A unique cluster name. Must be 3-25 lowercase alphanumeric characters or hyphens, start/end with alphanumeric, and contain no consecutive hyphens. Forces replacement on change.

- `network_id` **required** *string* &rarr;  The ID of the VPC network where the cluster will be deployed. The network must use the `sprut` SDN. **Forces replacement** on change.

- `network_plugin` **required** *string* &rarr;  The CNI plugin to use. Currently only `calico` is supported. **Forces replacement** on change.

- `pods_ipv4_cidr` **required** *string* &rarr;  The IPv4 CIDR block for the pod network. Must be a valid private CIDR (e.g., `10.100.0.0/16`). **Forces replacement** on change.

- `subnet_id` **required** *string* &rarr;  The ID of the subnet for cluster nodes. **Forces replacement** on change.

- `version` **required** *string* &rarr;  The Kubernetes version to deploy. Only upgrades to a higher version are allowed.

- `description` optional *string* &rarr;  A human-readable description of the cluster. **Forces replacement** on change. The maximum length is 256 characters.

- `external_network_id` optional *string* &rarr;  The ID of the external network for internet access. If omitted, the system selects one automatically. **Forces replacement** on change.

- `insecure_registries` optional *set of* *string* &rarr;  A set of registry addresses (e.g., `myregistry.com`) that can be accessed without TLS verification. **Forces replacement** on change.

- `labels` optional *map of* *string* &rarr;  A map of Kubernetes labels to apply to the cluster. Keys and values must conform to Kubernetes label syntax. **Forces replacement** on change.

- `loadbalancer_allowed_cidrs` optional *set of* *string* &rarr;  A set of CIDR blocks permitted to access the API load balancer. If empty, all IPs are allowed. **Forces replacement** on change.

- `public_ip` optional *boolean* &rarr;  If `true`, a floating IP is assigned to the API endpoint. `Forces replacement` on change. Default is `false`.

- `region` optional *string* &rarr;  The region in which to create the cluster. If omitted, the provider's `region` is used. **Forces replacement** on change.

- `uuid` optional *string* &rarr;  The cluster UUID. Can be specified; if not, it is generated automatically. **Changing this forces a new cluster.**


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `api_address` *string* &rarr;  The URL of the Kubernetes API server.

- `api_lb_fip` *string* &rarr;  The floating IP of the API load balancer (present if `public_ip` = `true`).

- `api_lb_vip` *string* &rarr;  The internal VIP of the API load balancer.

- `created_at` *string* &rarr;  The timestamp when the cluster was created (ISO 8601 format).

- `id` *string* &rarr;  The unique identifier of the cluster.

- `k8s_config` *string* &rarr;  The kubeconfig file contents. **Sensitive** — store securely.

- `master_disks`  *set* &rarr;  Information about the master node disks. Each disk includes `type` and `size`.
    - `size` *number* &rarr;  The size of the disk in GB.

    - `type` *string* &rarr;  The storage volume type provisioned for this master disk component.


- `node_groups`  *set* &rarr;  A list of node groups associated with the cluster. Each object contains `id`, `name`, `flavor`, `node_count`, and `availability_zone`.
    - `availability_zone` *string* &rarr;  The AZ where the node group resides.

    - `flavor` *string* &rarr;  The flavor ID used for the worker nodes.

    - `id` *string* &rarr;  The ID of the node group.

    - `name` *string* &rarr;  The name of the node group.

    - `node_count` *number* &rarr;  The current number of nodes in the group.


- `project_id` *string* &rarr;  The project ID that owns the cluster.

- `status` *string* &rarr;  The current cluster status (e.g., `RUNNING`, `RECONCILING`).



## Import

Clusters can be imported using the `id`, e.g.

```shell
terraform import vkcs_kubernetes_cluster_v2.k8s_cluster 39UdIv4W0EegBs2EYVeGdas38do
```
