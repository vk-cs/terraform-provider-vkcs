---
layout: "vkcs"
page_title: "vkcs: kubernetes_cluster"
description: |-
  Manages a kubernetes cluster.
---

# vkcs\_kubernetes\_cluster

Provides a kubernetes cluster resource. This can be used to create, modify and delete kubernetes clusters.

## Example Usage

```terraform

resource "vkcs_kubernetes_cluster" "mycluster" {
      name                = "terracluster"
      cluster_template_id = example_template_id
      master_flavor       = example_flavor_id
      master_count        = 1
      network_id          = example_network_id
      subnet_id           = example_subnet_id
      availability_zone   = "MS1"
      labels = {
        ingress_controller="nginx"
      }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the cluster. Changing this creates a new cluster. Should match the pattern `^[a-zA-Z][a-zA-Z0-9_.-]*$`.

* `cluster_template_id` - (Required) The UUID of the Kubernetes cluster
    template. It can be obtained using the cluster_template data source.

* `master_flavor` - (Optional) The UUID of a flavor for the master nodes.
 If master_flavor is not present, value from cluster_template will be used.

* `network_id` - (Required) The UUID of the network that will be attached to the cluster.
 Changing this creates a new cluster.

* `subnet_id` - (Required) The UUID of the subnet that will be attached to the cluster.
 Changing this creates a new cluster.

* `keypair` - (Optional) The name of the Compute service SSH keypair. Changing
    this creates a new cluster.

* `labels` - (Optional) The list of optional key value pairs representing additional
    properties of the cluster. Changing this creates a new cluster.
  * `docker_registry_enabled=true` to preinstall Docker Registry.
  * `prometheus_monitoring=true` to preinstall monitoring system based on Prometheus and Grafana.
  * `ingress_controller="nginx"` to preinstall NGINX Ingress Controller.

* `master_count` - (Optional) The number of master nodes for the cluster.
    Changing this creates a new cluster.
    
* `pods_network_cidr` - (Optional) The network cidr used in k8s virtual network.

* `floating_ip_enabled` - (Required) Floating ip is enabled.

* `api_lb_vip` - (Optional) API LoadBalancer vip.

* `api_lb_fip` - (Optional) API LoadBalancer fip.

* `registry_auth_password` - (Optional) Docker registry access password.

* `availability_zone` - (Required) Zones available for cluster. `GZ1` and `MS1` zones are available. **New since v0.3.3**.

* `region` - (Optional) Region to use for the cluster. Default is a region configured for provider. **New since v0.4.0**.

* `loadbalancer_subnet_id` - (Optional) The UUID of the load balancer's subnet. Changing this creates new cluster. **New since v0.5.4**.

## Attributes

This resource exports the following attributes:

* `name` - The name of the cluster.
* `project_id` - The project of the cluster.
* `created_at` - The time at which cluster was created.
* `updated_at` - The time at which cluster was created.
* `api_address` - COE API address.
* `cluster_template_id` - The UUID of the V1 Container Infra cluster template.
* `discovery_url` - The URL used for cluster node discovery.
* `master_flavor` - The UUID of a flavor for the master nodes. 
* `keypair` - The name of the Compute service SSH keypair.
* `labels` - The list of key value pairs representing additional properties of
                 the cluster.
* `master_count` - The number of master nodes for the cluster.
* `master_addresses` - IP addresses of the master node of the cluster.
* `stack_id` - UUID of the Orchestration service stack.
* `network_id` - UUID of the cluster's network.
* `subnet_id` - UUID of the cluster's subnet.
* `status` - Current state of a cluster. Changing this to `RUNNING` or `SHUTOFF` will turn cluster on/off.
* `pods_network_cidr` - Network cidr of k8s virtual network
* `floating_ip_enabled` - Floating ip is enabled.
* `api_lb_vip` - API LoadBalancer vip.
* `api_lb_fip` - API LoadBalancer fip.
* `ingress_floating_ip` - Floating IP created for ingress service.
* `registry_auth_password` - Docker registry access password.
* `availability_zone` - Availability zone of the cluster. **New since v0.3.3**
* `loadbalancer_subnet_id` - UUID of the load balancer's subnet. **New since v0.5.4**.

## Import

Clusters can be imported using the `id`, e.g.

```
$ terraform import vkcs_kubernetes_cluster.mycluster ce0f9463-dd25-474b-9fe8-94de63e5e42b
```