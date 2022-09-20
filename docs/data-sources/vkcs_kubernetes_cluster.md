---
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_cluster"
description: |-
  Get information on cluster.
---

# vkcs_kubernetes_cluster

Use this data source to get the ID of an available VKCS kubernetes cluster.

## Example Usage

```terraform
data "vkcs_kubernetes_cluster" "mycluster" {
  name = "myclustername"
}
```
```terraform
data "vkcs_kubernetes_cluster" "mycluster" {
  cluster_id = "myclusteruuid"
}
```
## Argument Reference
- `cluster_id` **String** (*Optional*) The UUID of the Kubernetes cluster template. **Note**: Only one of `name` or `cluster_id` must be specified.

- `dns_domain` **String** (*Optional*) Custom DNS cluster domain.

- `name` **String** (*Optional*) The name of the cluster. **Note**: Only one of `name` or `cluster_id` must be specified.

- `region` **String** (*Optional*) The region in which to obtain the Container Infra client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
- `cluster_id` **String** See Argument Reference above.

- `dns_domain` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `api_address` **String** COE API address.

- `api_lb_fip` **String** API LoadBalancer fip.

- `api_lb_vip` **String** API LoadBalancer vip.

- `availability_zone` **String** Availability zone of the cluster.

- `cluster_template_id` **String** The UUID of the V1 Container Infra cluster template.

- `created_at` **String** The time at which cluster was created.

- `discovery_url` **String** The URL used for cluster node discovery.

- `floating_ip_enabled` **Boolean** Indicates whether floating ip is enabled for cluster.

- `id` **String** ID of the resource.

- `ingress_floating_ip` **String** Floating IP created for ingress service.

- `insecure_registries` **String** Addresses of registries from which you can download images without checking certificates.

- `k8s_config` **String** Kubeconfig for cluster

- `keypair` **String** The name of the Compute service SSH keypair.

- `labels` <strong>Map of </strong>**String** The list of key value pairs representing additional properties of the cluster.

- `loadbalancer_subnet_id` **String** The ID of load balancer's subnet.

- `master_addresses` **String** IP addresses of the master node of the cluster.

- `master_count` **Number** The number of master nodes for the cluster.

- `master_flavor` **String** The ID of the flavor for the master nodes.

- `network_id` **String** UUID of the cluster's network.

- `pods_network_cidr` **String** Network cidr of k8s virtual network.

- `project_id` **String** The project of the cluster.

- `registry_auth_password` **String** Docker registry access password.

- `stack_id` **String** UUID of the Orchestration service stack.

- `status` **String** Current state of a cluster.

- `subnet_id` **String** UUID of the cluster's subnet.

- `updated_at` **String** The time at which cluster was created.

- `user_id` **String** The user of the cluster.


