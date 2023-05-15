---
subcategory: "Kubernetes"
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
  version = "1.24"
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
- `availability_zone` **required** *string* &rarr;  Availability zone of the cluster.

- `cluster_template_id` **required** *string* &rarr;  The UUID of the Kubernetes cluster template. It can be obtained using the cluster_template data source.

- `floating_ip_enabled` **required** *boolean* &rarr;  Floating ip is enabled.

- `name` **required** *string* &rarr;  The name of the cluster. Changing this creates a new cluster. Should match the pattern `^[a-zA-Z][a-zA-Z0-9_.-]*$`.

- `network_id` **required** *string* &rarr;  UUID of the cluster's network.

- `subnet_id` **required** *string* &rarr;  UUID of the cluster's subnet.

- `api_lb_fip` optional *string* &rarr;  API LoadBalancer fip.

- `api_lb_vip` optional *string* &rarr;  API LoadBalancer vip.

- `dns_domain` optional *string* &rarr;  Custom DNS cluster domain. Changing this creates a new cluster.

- `ingress_floating_ip` optional *string* &rarr;  Floating IP created for ingress service.

- `insecure_registries` optional *string* &rarr;  Addresses of registries from which you can download images without checking certificates. Changing this creates a new cluster.

- `keypair` optional *string* &rarr;  The name of the Compute service SSH keypair. Changing this creates a new cluster.

- `labels` optional *map of* *string* &rarr;  The list of optional key value pairs representing additional properties of the cluster. Changing this creates a new cluster.

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

- `loadbalancer_subnet_id` optional *string* &rarr;  The UUID of the load balancer's subnet. Changing this creates new cluster.

- `master_count` optional *number* &rarr;  The number of master nodes for the cluster. Changing this creates a new cluster.

- `master_flavor` optional *string* &rarr;  The UUID of a flavor for the master nodes. If master_flavor is not present, value from cluster_template will be used.

- `pods_network_cidr` optional *string* &rarr;  Network cidr of k8s virtual network

- `region` optional *string* &rarr;  Region to use for the cluster. Default is a region configured for provider.

- `registry_auth_password` optional *string* &rarr;  Docker registry access password.

- `status` optional *string* &rarr;  Current state of a cluster. Changing this to `RUNNING` or `SHUTOFF` will turn cluster on/off.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `api_address` *string* &rarr;  COE API address.

- `created_at` *string* &rarr;  The time at which cluster was created.

- `id` *string* &rarr;  ID of the resource.

- `master_addresses` *string* &rarr;  IP addresses of the master node of the cluster.

- `project_id` *string* &rarr;  The project of the cluster.

- `stack_id` *string* &rarr;  UUID of the Orchestration service stack.

- `updated_at` *string* &rarr;  The time at which cluster was created.

- `user_id` *string* &rarr;  The user of the cluster.



## Import

Clusters can be imported using the `id`, e.g.

```shell
terraform import vkcs_kubernetes_cluster.mycluster ce0f9463-dd25-474b-9fe8-94de63e5e42b
```
