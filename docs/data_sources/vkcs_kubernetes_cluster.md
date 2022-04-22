---
layout: "vkcs"
page_title: "vkcs: kubernetes_cluster"
description: |-
  Get information on cluster.
---

# vkcs\_kubernetes\_cluster

Use this data source to get the ID of an available VKCS kubernetes cluster.

## Example Usage
```hcl
data "vkcs_kubernetes_cluster" "mycluster" {
  name = "myclustername"
}
```
```hcl
data "vkcs_kubernetes_cluster" "mycluster" {
  cluster_id = "myclusteruuid"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the cluster.

* `cluster_id` - (Optional) The UUID of the Kubernetes cluster
    template.

* `region` - (Optional) The region in which to obtain the Container Infra
    client.
    If omitted, the `region` argument of the provider is used.
        
**Note**: Only one of `name` or `cluster_id` must be specified

    
## Attributes
`id` is set to the ID of the found cluster template. In addition, the following
attributes are exported:

* `api_address` - COE API address.
* `cluster_template_id` - The UUID of the V1 Container Infra cluster template.
* `created_at` - The time at which cluster was created.
* `discovery_url` - The URL used for cluster node discovery.
* `k8s_config` - Kubeconfig for cluster
* `keypair` - The name of the Compute service SSH keypair.
* `labels` - The list of key value pairs representing additional properties of the cluster.
* `loadbalancer_subnet_id` - The ID of load balancer's subnet.
* `master_addresses` - IP addresses of the master node of the cluster.
* `master_count` - The number of master nodes for the cluster.
* `master_flavor` - The ID of the flavor for the master nodes.
* `name` - The name of the cluster.
* `network_id` - UUID of the cluster's network.
* `pods_network_cidr` - Network cidr of k8s virtual network
* `project_id` - The project of the cluster.
* `stack_id` - UUID of the Orchestration service stack.
* `status` - Current state of a cluster.
* `subnet_id` - UUID of the cluster's subnet.
* `updated_at` - The time at which cluster was created.
* `insecure_registries` - Addresses of registries from which you can download images without checking certificates.
