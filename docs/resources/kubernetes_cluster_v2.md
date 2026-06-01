---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_cluster_v2"
description: |-
  Provides a Kubernetes cluster resource. This can be used to create, modify and delete Kubernetes clusters.
---

# vkcs_kubernetes_cluster_v2

Provides a Kubernetes cluster resource. This can be used to create, modify and delete Kubernetes clusters.

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
- `availability_zones` **required** *set of* *string* &rarr;  A set of availability zones for the cluster masters. A `standard` cluster must specify one zone, while a `regional` cluster must specify three zones. **Forces replacement** on change.

- `cluster_type` **required** *string* &rarr;  Either `standard` or `regional`. A `standard` cluster uses a single availability zone, while a `regional` cluster uses three availability zones. **Forces replacement** on change.

- `loadbalancer_subnet_id` **required** *string* &rarr;  The ID of the subnet used for allocating IP addresses for LoadBalancer services. **Forces replacement** on change.

- `master_count` **required** *number* &rarr;  The number of master nodes. Must be an odd number (1, 3, or 5) to maintain etcd quorum. **Forces replacement** on change.

- `master_flavor` **required** *string* &rarr;  The flavor ID used for master nodes. Changing this triggers a rolling upgrade of the control plane.

- `name` **required** *string* &rarr;  A unique cluster name. Must be 3-25 lowercase alphanumeric characters or hyphens, start and end with an alphanumeric character, and must not contain consecutive hyphens. **Forces replacement** on change.

- `network_id` **required** *string* &rarr;  The ID of the network in which the cluster will be deployed. The network must use the `sprut` SDN. **Forces replacement** on change.

- `network_plugin` **required** *string* &rarr;  The CNI plugin to use. Currently, only `calico` is supported. **Forces replacement** on change.

- `pods_ipv4_cidr` **required** *string* &rarr;  The IPv4 CIDR block used for the pod network. Must be a valid private CIDR (e.g., `10.100.0.0/16`). **Forces replacement** on change.

- `subnet_id` **required** *string* &rarr;  The ID of the subnet in which the cluster worker nodes will be deployed. **Forces replacement** on change.

- `version` **required** *string* &rarr;  The Kubernetes version to deploy. Only upgrades to a higher version are allowed.

- `description` optional *string* &rarr;  A human-readable description of the cluster. The maximum length is 256 characters. **Forces replacement** on change.

- `external_network_id` optional *string* &rarr;  The ID of the external network used for internet access. If omitted, one is selected automatically. **Forces replacement** on change.

- `insecure_registries` optional *set of* *string* &rarr;  A set of container registry hosts (e.g., `myregistry.com`) that can be accessed without TLS verification. **Forces replacement** on change.

- `labels` optional *map of* *string* &rarr;  A map of Kubernetes labels applied to the cluster. Keys and values must conform to Kubernetes label syntax. **Forces replacement** on change.

- `loadbalancer_allowed_cidrs` optional *set of* *string* &rarr;  A set of CIDR blocks allowed to access the API load balancer. If empty, access is allowed from all IP addresses. **Forces replacement** on change.

- `public_ip` optional *boolean* &rarr;  If `true`, assigns a floating IP to the API endpoint. When enabled, `external_network_id` must be specified. **Forces replacement** on change. Default is `false`.

- `region` optional *string* &rarr;  The region in which to create the cluster. If omitted, the provider's `region` is used. **Forces replacement** on change.

- `uuid` optional *string* &rarr;  The UUID of the cluster. If omitted, it is generated automatically. **Forces replacement** on change.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `api_address` *string* &rarr;  The Kubernetes API server URL.

- `api_lb_fip` *string* &rarr;  The floating IP of the API load balancer, present only if `public_ip` is `true`.

- `api_lb_vip` *string* &rarr;  The internal virtual IP (VIP) of the API load balancer.

- `created_at` *string* &rarr;  The time when the cluster was created.

- `id` *string* &rarr;  The unique identifier of the cluster.

- `k8s_config` *string* &rarr;  The contents of the kubeconfig file. **Sensitive** — store securely.

- `master_disks`  *set* &rarr;  Information about the master node disks.
    - `size` *number* &rarr;  The size of the disk in gigabytes (GB).

    - `type` *string* &rarr;  The storage volume type for this master disk.


- `node_groups`  *set* &rarr;  A list of node groups associated with the cluster.
    - `availability_zone` *string* &rarr;  The availability zone where the node group is located.

    - `flavor` *string* &rarr;  The flavor ID used for worker nodes.

    - `id` *string* &rarr;  The ID of the node group.

    - `name` *string* &rarr;  The name of the node group.

    - `node_count` *number* &rarr;  The current number of nodes in the node group.


- `project_id` *string* &rarr;  The ID of the project that owns the cluster.

- `status` *string* &rarr;  The current status of the cluster (e.g., `RUNNING`, `RECONCILING`).



## Import

Clusters can be imported using the `id`, e.g.

```shell
terraform import vkcs_kubernetes_cluster_v2.k8s_cluster 39UdIv4W0EegBs2EYVeGdas38do
```
