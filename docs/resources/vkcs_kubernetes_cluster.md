---
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_cluster"
description: |-
  Manages a kubernetes cluster.
---

# vkcs_kubernetes_cluster

Provides a kubernetes cluster resource. This can be used to create, modify and delete kubernetes clusters.

## Example Usage
```terraform
data "vkcs_kubernetes_clustertemplate" "ct" {
  version = "1.21.4"
}

resource "vkcs_kubernetes_cluster" "k8s-cluster" {
  depends_on = [
    vkcs_networking_router_interface.k8s,
  ]

  name                = "k8s-cluster"
  cluster_template_id = data.vkcs_kubernetes_clustertemplate.ct.id
  master_flavor       = data.vkcs_compute_flavor.k8s.id
  master_count        = 1

  network_id          = vkcs_networking_network.k8s.id
  subnet_id           = vkcs_networking_subnet.k8s-subnetwork.id
  floating_ip_enabled = true
  availability_zone   = "MS1"
  insecure_registries = ["1.2.3.4"]
}
```
## Argument Reference
- `availability_zone` **String** (***Required***) Availability zone of the cluster.

- `cluster_template_id` **String** (***Required***) The UUID of the Kubernetes cluster template. It can be obtained using the cluster_template data source.

- `floating_ip_enabled` **Boolean** (***Required***) Floating ip is enabled.

- `name` **String** (***Required***) The name of the cluster. Changing this creates a new cluster. Should match the pattern `^[a-zA-Z][a-zA-Z0-9_.-]*$`.

- `network_id` **String** (***Required***) UUID of the cluster's network.

- `subnet_id` **String** (***Required***) UUID of the cluster's subnet.

- `api_lb_fip` **String** (*Optional*) API LoadBalancer fip.

- `api_lb_vip` **String** (*Optional*) API LoadBalancer vip.

- `dns_domain` **String** (*Optional*) Custom DNS cluster domain. Changing this creates a new cluster.

- `ingress_floating_ip` **String** (*Optional*) Floating IP created for ingress service.

- `insecure_registries` **String** (*Optional*) Addresses of registries from which you can download images without checking certificates. Changing this creates a new cluster.

- `keypair` **String** (*Optional*) The name of the Compute service SSH keypair. Changing this creates a new cluster.

- `labels` <strong>Map of </strong>**String** (*Optional*) The list of optional key value pairs representing additional properties of the cluster. Changing this creates a new cluster.

  * `calico_ipv4pool` to set subnet where pods will be created. Default 10.100.0.0/16.
  * `clean_volumes` to remove pvc volumes when deleting a cluster. Default False.
  * `cloud_monitoring` to enable cloud monitoring feature.
  * `docker_registry_enabled=true` to preinstall Docker Registry.
  * `etcd_volume_size` to set etcd volume size. Default 10Gb.
  * `ingress_controller="nginx"` to preinstall NGINX Ingress Controller.
  * `kube_log_level` to set log level for kubelet in range 0 to 8.
  * `master_volume_size` to set master vm volume size. Default 50Gb.
  * `cluster_node_volume_type` to set master vm volume type. Default ceph-hdd.
  * `prometheus_monitoring=true` to preinstall monitoring system based on Prometheus and Grafana.

- `loadbalancer_subnet_id` **String** (*Optional*) The UUID of the load balancer's subnet. Changing this creates new cluster.

- `master_count` **Number** (*Optional*) The number of master nodes for the cluster. Changing this creates a new cluster.

- `master_flavor` **String** (*Optional*) The UUID of a flavor for the master nodes. If master_flavor is not present, value from cluster_template will be used.

- `pods_network_cidr` **String** (*Optional*) Network cidr of k8s virtual network

- `region` **String** (*Optional*) Region to use for the cluster. Default is a region configured for provider.

- `registry_auth_password` **String** (*Optional*) Docker registry access password.

- `status` **String** (*Optional*) Current state of a cluster. Changing this to `RUNNING` or `SHUTOFF` will turn cluster on/off.


## Attributes Reference
- `availability_zone` **String** See Argument Reference above.

- `cluster_template_id` **String** See Argument Reference above.

- `floating_ip_enabled` **Boolean** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `network_id` **String** See Argument Reference above.

- `subnet_id` **String** See Argument Reference above.

- `api_lb_fip` **String** See Argument Reference above.

- `api_lb_vip` **String** See Argument Reference above.

- `dns_domain` **String** See Argument Reference above.

- `ingress_floating_ip` **String** See Argument Reference above.

- `insecure_registries` **String** See Argument Reference above.

- `keypair` **String** See Argument Reference above.

- `labels` <strong>Map of </strong>**String** See Argument Reference above.

- `loadbalancer_subnet_id` **String** See Argument Reference above.

- `master_count` **Number** See Argument Reference above.

- `master_flavor` **String** See Argument Reference above.

- `pods_network_cidr` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `registry_auth_password` **String** See Argument Reference above.

- `status` **String** See Argument Reference above.

- `api_address` **String** COE API address.

- `created_at` **String** The time at which cluster was created.

- `id` **String** ID of the resource.

- `master_addresses` **String** IP addresses of the master node of the cluster.

- `project_id` **String** The project of the cluster.

- `stack_id` **String** UUID of the Orchestration service stack.

- `updated_at` **String** The time at which cluster was created.

- `user_id` **String** The user of the cluster.



## Import

Clusters can be imported using the `id`, e.g.

```shell
terraform import vkcs_kubernetes_cluster.mycluster ce0f9463-dd25-474b-9fe8-94de63e5e42b
```
