---
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
- `cluster_template_uuid` **String** (*Optional*) The UUID of the cluster template. **Note**: Only one of `name` or `version` or `cluster_template_uuid` must be specified.

- `name` **String** (*Optional*) The name of the cluster template. **Note**: Only one of `name` or `version` or `cluster_template_uuid` must be specified.

- `region` **String** (*Optional*) The region in which to obtain the V1 Container Infra client. If omitted, the `region` argument of the provider is used.

- `version` **String** (*Optional*) Kubernetes version of the cluster.


## Attributes Reference
- `cluster_template_uuid` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `version` **String** See Argument Reference above.

- `apiserver_port` **Number** The API server port for the Container Orchestration Engine for this cluster template.

- `cluster_distro` **String** The distro for the cluster (fedora-atomic, coreos, etc.).

- `created_at` **String** The time at which cluster template was created.

- `deprecated_at` **String** The time at which the cluster template is deprecated.

- `dns_nameserver` **String** Address of the DNS nameserver that is used in nodes of the cluster.

- `docker_storage_driver` **String** Docker storage driver. Changing this updates the Docker storage driver of the existing cluster template.

- `docker_volume_size` **Number** The size (in GB) of the Docker volume.

- `external_network_id` **String** The ID of the external network that will be used for the cluster.

- `flavor` **String** The ID of flavor for the nodes of the cluster.

- `floating_ip_enabled` **Boolean** Indicates whether created cluster should create IP floating IP for every node or not.

- `id` **String** ID of the resource.

- `image` **String** The reference to an image that is used for nodes of the cluster.

- `insecure_registry` **String** The insecure registry URL for the cluster template.

- `keypair_id` **String** The name of the Compute service SSH keypair.

- `labels` <strong>Map of </strong>**String** The list of key value pairs representing additional properties of the cluster template.

- `master_flavor` **String** The ID of flavor for the master nodes.

- `master_lb_enabled` **Boolean** Indicates whether created cluster should has a loadbalancer for master nodes or not.

- `network_driver` **String** The name of the driver for the container network.

- `no_proxy` **String** A comma-separated list of IP addresses that shouldn't be used in the cluster.

- `project_id` **String** The project of the cluster template.

- `public` **Boolean** Indicates whether cluster template should be public.

- `registry_enabled` **Boolean** Indicates whether Docker registry is enabled in the cluster.

- `server_type` **String** The server type for the cluster template.

- `tls_disabled` **Boolean** Indicates whether the TLS should be disabled in the cluster.

- `updated_at` **String** The time at which cluster template was updated.

- `user_id` **String** The user of the cluster template.

- `volume_driver` **String** The name of the driver that is used for the volumes of the cluster nodes.


