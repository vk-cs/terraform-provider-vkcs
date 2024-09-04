---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_cluster"
description: |-
  Get information on cluster.
---

# vkcs_kubernetes_cluster

Use this data source to get the ID of an available VKCS kubernetes cluster.

## Example Usage
```terraform
data "vkcs_kubernetes_cluster" "k8s-cluster" {
  cluster_id = vkcs_kubernetes_cluster.k8s-cluster.id
}
```

## Argument Reference
- `cluster_id` optional *string* &rarr;  The UUID of the Kubernetes cluster template. <br>**Note:** Only one of `name` or `cluster_id` must be specified.

- `dns_domain` optional *string* &rarr;  Custom DNS cluster domain.

- `name` optional *string* &rarr;  The name of the cluster. <br>**Note:** Only one of `name` or `cluster_id` must be specified.

- `region` optional *string* &rarr;  The region in which to obtain the Container Infra client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `api_address` *string* &rarr;  COE API address.

- `api_lb_fip` *string* &rarr;  API LoadBalancer fip.

- `api_lb_vip` *string* &rarr;  API LoadBalancer vip.

- `availability_zone` *string* &rarr;  Availability zone of the cluster.

- `availability_zones` *set of* *string* &rarr;  Availability zones of the regional cluster<br>**New since v0.8.3**.

- `cluster_template_id` *string* &rarr;  The UUID of the V1 Container Infra cluster template.

- `cluster_type` *string* &rarr;  Type of the kubernetes cluster, may be `standard` or `regional`<br>**New since v0.8.3**.

- `created_at` *string* &rarr;  The time at which cluster was created.

- `discovery_url` *string* &rarr;  The URL used for cluster node discovery.

- `floating_ip_enabled` *boolean* &rarr;  Indicates whether floating ip is enabled for cluster.

- `id` *string* &rarr;  ID of the resource.

- `ingress_floating_ip` *string* &rarr;  Floating IP created for ingress service.

- `insecure_registries` *string* &rarr;  Addresses of registries from which you can download images without checking certificates.

- `k8s_config` *string* &rarr;  Kubeconfig for cluster

- `keypair` *string* &rarr;  The name of the Compute service SSH keypair.

- `labels` *map of* *string* &rarr;  The list of key value pairs representing additional properties of the cluster.

- `loadbalancer_subnet_id` *string* &rarr;  The ID of load balancer's subnet.

- `master_addresses` *string* &rarr;  IP addresses of the master node of the cluster.

- `master_count` *number* &rarr;  The number of master nodes for the cluster.

- `master_flavor` *string* &rarr;  The ID of the flavor for the master nodes.

- `network_id` *string* &rarr;  UUID of the cluster's network.

- `pods_network_cidr` *string* &rarr;  Network cidr of k8s virtual network.

- `project_id` *string* &rarr;  The project of the cluster.

- `registry_auth_password` *string* &rarr;  Docker registry access password.

- `stack_id` *string* &rarr;  UUID of the Orchestration service stack.

- `status` *string* &rarr;  Current state of a cluster.

- `subnet_id` *string* &rarr;  UUID of the cluster's subnet.

- `sync_security_policy` *boolean* &rarr;  Enables syncing of security policies of cluster.<br>**New since v0.7.0**.

- `updated_at` *string* &rarr;  The time at which cluster was created.

- `user_id` *string* &rarr;  The user of the cluster.


