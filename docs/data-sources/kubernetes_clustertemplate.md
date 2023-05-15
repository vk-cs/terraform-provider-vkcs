---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_clustertemplate"
description: |-
  Get information on an VKCS kubernetes cluster template.
---

# vkcs_kubernetes_clustertemplate

Use this data source to get the ID of an available VKCS kubernetes cluster template.

## Example Usage

```terraform
data "vkcs_kubernetes_clustertemplate" "example_template" {
  name = "clustertemplate_1"
}

output "example_template_id" {
  value = "${data.vkcs_kubernetes_clustertemplate.example_template.id}"
}
```
```terraform
data "vkcs_kubernetes_clustertemplate" "example_template_by_version" {
  version = "1.20.4"
}

output "example_template_id" {
  value = "${data.vkcs_kubernetes_clustertemplate.example_template_by_version.id}"
}
```
## Argument Reference
- `cluster_template_uuid` optional *string* &rarr;  The UUID of the cluster template. **Note**: Only one of `name` or `version` or `cluster_template_uuid` must be specified.

- `name` optional *string* &rarr;  The name of the cluster template. **Note**: Only one of `name` or `version` or `cluster_template_uuid` must be specified.

- `region` optional *string* &rarr;  The region in which to obtain the V1 Container Infra client. If omitted, the `region` argument of the provider is used.

- `version` optional *string* &rarr;  Kubernetes version of the cluster.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `apiserver_port` *number* &rarr;  The API server port for the Container Orchestration Engine for this cluster template.

- `cluster_distro` *string* &rarr;  The distro for the cluster (fedora-atomic, coreos, etc.).

- `created_at` *string* &rarr;  The time at which cluster template was created.

- `deprecated_at` *string* &rarr;  The time at which the cluster template is deprecated.

- `dns_nameserver` *string* &rarr;  Address of the DNS nameserver that is used in nodes of the cluster.

- `docker_storage_driver` *string* &rarr;  Docker storage driver. Changing this updates the Docker storage driver of the existing cluster template.

- `docker_volume_size` *number* &rarr;  The size (in GB) of the Docker volume.

- `external_network_id` *string* &rarr;  The ID of the external network that will be used for the cluster.

- `flavor` *string* &rarr;  The ID of flavor for the nodes of the cluster.

- `floating_ip_enabled` *boolean* &rarr;  Indicates whether created cluster should create IP floating IP for every node or not.

- `id` *string* &rarr;  ID of the resource.

- `image` *string* &rarr;  The reference to an image that is used for nodes of the cluster.

- `insecure_registry` *string* &rarr;  The insecure registry URL for the cluster template.

- `keypair_id` *string* &rarr;  The name of the Compute service SSH keypair.

- `labels` *map of* *string* &rarr;  The list of key value pairs representing additional properties of the cluster template.

- `master_flavor` *string* &rarr;  The ID of flavor for the master nodes.

- `master_lb_enabled` *boolean* &rarr;  Indicates whether created cluster should has a loadbalancer for master nodes or not.

- `network_driver` *string* &rarr;  The name of the driver for the container network.

- `no_proxy` *string* &rarr;  A comma-separated list of IP addresses that shouldn't be used in the cluster.

- `project_id` *string* &rarr;  The project of the cluster template.

- `public` *boolean* &rarr;  Indicates whether cluster template should be public.

- `registry_enabled` *boolean* &rarr;  Indicates whether Docker registry is enabled in the cluster.

- `server_type` *string* &rarr;  The server type for the cluster template.

- `tls_disabled` *boolean* &rarr;  Indicates whether the TLS should be disabled in the cluster.

- `updated_at` *string* &rarr;  The time at which cluster template was updated.

- `user_id` *string* &rarr;  The user of the cluster template.

- `volume_driver` *string* &rarr;  The name of the driver that is used for the volumes of the cluster nodes.


