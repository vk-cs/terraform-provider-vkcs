---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_cluster"
description: |-
  Manages a kubernetes cluster.
---

# vkcs_kubernetes_cluster

Provides a kubernetes cluster resource. This can be used to create, modify and delete kubernetes clusters.

## Standard Kubernetes cluster
```terraform
resource "vkcs_kubernetes_cluster" "k8s-cluster" {
  name                = "k8s-standard-cluster"
  cluster_type        = "standard"
  cluster_template_id = data.vkcs_kubernetes_clustertemplate.k8s_24.id
  master_flavor       = data.vkcs_compute_flavor.basic.id
  master_count        = 1

  labels = {
    cloud_monitoring         = "true"
    kube_log_level           = "2"
    clean_volumes            = "true"
    master_volume_size       = "100"
    cluster_node_volume_type = "ceph-ssd"
  }

  availability_zone   = "MS1"
  network_id          = vkcs_networking_network.app.id
  subnet_id           = vkcs_networking_subnet.app.id
  floating_ip_enabled = true

  sync_security_policy = true
  # If your configuration also defines a network for the instance,
  # ensure it is attached to a router before creating of the instance
  depends_on = [
    vkcs_networking_router_interface.app,
  ]
}
```

## Regional Kubernetes cluster
```terraform
resource "vkcs_kubernetes_cluster" "k8s-cluster" {
  name                = "k8s-regional-cluster"
  cluster_type        = "regional"
  cluster_template_id = data.vkcs_kubernetes_clustertemplate.k8s_24.id
  master_flavor       = data.vkcs_compute_flavor.basic.id
  master_count        = 3

  labels = {
    cloud_monitoring         = "true"
    kube_log_level           = "2"
    clean_volumes            = "true"
    master_volume_size       = "100"
    cluster_node_volume_type = "ceph-ssd"
  }

  network_id          = vkcs_networking_network.app.id
  subnet_id           = vkcs_networking_subnet.app.id
  floating_ip_enabled = true

  sync_security_policy = true
  # If your configuration also defines a network for the instance,
  # ensure it is attached to a router before creating of the instance
  depends_on = [
    vkcs_networking_router_interface.app,
  ]
}
```

## Argument Reference
- `cluster_template_id` **required** *string* &rarr;  The UUID of the Kubernetes cluster template. It can be obtained using the cluster_template data source.

- `floating_ip_enabled` **required** *boolean* &rarr;  Floating ip is enabled.

- `name` **required** *string* &rarr;  The name of the cluster. Changing this creates a new cluster. Should match the pattern `^[a-zA-Z][a-zA-Z0-9_.-]*$`.

- `network_id` **required** *string* &rarr;  UUID of the cluster's network.

- `subnet_id` **required** *string* &rarr;  UUID of the cluster's subnet.

- `api_lb_fip` optional *string* &rarr;  API LoadBalancer fip. IP address field.

- `api_lb_vip` optional *string* &rarr;  API LoadBalancer vip. IP address field.

- `availability_zone` optional *string* &rarr;  Availability zone of the cluster, set this argument only for cluster with type `standard`.

- `availability_zones` optional *set of* *string* &rarr;  Availability zones of the regional cluster, set this argument only for cluster with type `regional`. If you do not set this argument, the availability zones will be selected automatically.<br>**New since v0.8.3**.

- `cluster_type` optional *string* &rarr;  Type of the kubernetes cluster, may be `standard` or `regional`. Default type is `standard`.<br>**New since v0.8.3**.

- `dns_domain` optional *string* &rarr;  Custom DNS cluster domain. Changing this creates a new cluster.

- `ingress_floating_ip` optional *string* &rarr;  Floating IP created for ingress service.

- `insecure_registries` optional *string* &rarr;  Addresses of registries from which you can download images without checking certificates. Changing this creates a new cluster.

- `keypair` optional *string* &rarr;  The name of the Compute service SSH keypair. Changing this creates a new cluster.

- `labels` optional *map of* *string* &rarr;  The list of optional key value pairs representing additional properties of the cluster. <br>**Note:** Updating this attribute will not immediately apply the changes; these options will be used when recreating or deleting cluster nodes, for example, during an upgrade operation.

  * `calico_ipv4pool` to set subnet where pods will be created. Default 10.100.0.0/16. <br>**Note:** Updating this value while the cluster is running is dangerous because it can lead to loss of connectivity of the cluster nodes.
  * `clean_volumes` to remove pvc volumes when deleting a cluster. Default False. <br>**Note:** Changes to this value will be applied immediately.
  * `cloud_monitoring` to enable cloud monitoring feature. Default False.
  * `etcd_volume_size` to set etcd volume size in GB. Default 10.
  * `kube_log_level` to set log level for kubelet in range 0 to 8. Default 0.
  * `master_volume_size` to set master vm volume size in GB. Default 50.
  * `cluster_node_volume_type` to set master vm volume type. Default ceph-ssd.

- `loadbalancer_subnet_id` optional *string* &rarr;  The UUID of the load balancer's subnet. Changing this creates new cluster.

- `master_count` optional *number* &rarr;  The number of master nodes for the cluster. Changing this creates a new cluster.

- `master_flavor` optional *string* &rarr;  The UUID of a flavor for the master nodes. If master_flavor is not present, value from cluster_template will be used.

- `pods_network_cidr` optional *string* &rarr;  Network cidr of k8s virtual network

- `region` optional *string* &rarr;  Region to use for the cluster. Default is a region configured for provider.

- `registry_auth_password` optional *string* &rarr;  Docker registry access password.

- `status` optional *string* &rarr;  Current state of a cluster. Changing this to `RUNNING` or `SHUTOFF` will turn cluster on/off.

- `sync_security_policy` optional *boolean* &rarr;  Enables syncing of security policies of cluster. Default value is false.<br>**New since v0.7.0**.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `all_labels` *map of* *string* &rarr;  The read-only map of all cluster labels.<br>**New since v0.5.1**.

- `api_address` *string* &rarr;  COE API address.

- `created_at` *string* &rarr;  The time at which cluster was created.

- `id` *string* &rarr;  ID of the resource.

- `k8s_config` *string* &rarr;  Contents of the kubeconfig file. Use it to authenticate to Kubernetes cluster.<br>**New since v0.8.1**.

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
