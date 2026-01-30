---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_cluster_addon_v2"
description: |-
  Provides a Kubernetes cluster addon resource. This can be used to create, modify, and delete Kubernetes cluster addon.
---

# vkcs_kubernetes_cluster_addon_v2

Provides a Kubernetes cluster addon resource. This can be used to create, modify, and delete Kubernetes cluster addon.

## Example Usage
```terraform
resource "vkcs_kubernetes_cluster_addon_v2" "ingress_nginx" {
  cluster_id           = vkcs_kubernetes_cluster_v2.k8s_cluster.id
  addon_id             = data.vkcs_kubernetes_addon_v2.ingress_nginx.id
  namespace            = "ingress-nginx"
  configuration_values = data.vkcs_kubernetes_addon_v2.ingress_nginx.values_template

  depends_on = [
    vkcs_kubernetes_node_group_v2.k8s_node_group
  ]
}
```

## Argument Reference
- `addon_id` **required** *string* &rarr;  The unique identifier of the addon definition. **Forces replacement** on change.

- `addon_name` **required** *string* &rarr;  The human-readable name of the cluster addon. **Forces replacement** on change.

- `addon_version_id` **required** *string* &rarr;  The ID of the addon version to install. **Forces replacement** on change.

- `cluster_id` **required** *string* &rarr;  The ID of the Kubernetes cluster where the addon will be installed. **Forces replacement** on change.

- `namespace` **required** *string* &rarr;  The Kubernetes namespace where the addon will be deployed. Must be a valid DNS subdomain name (e.g., 'my-namespace'). **Forces replacement** on change.

- `values` **required** *string* &rarr;  The Helm chart values for configuring the addon, provided as YAML string.

- `region` optional *string* &rarr;  The region in which to create the cluster addon. If omitted, the provider's `region` is used. **Forces replacement** on change.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  The timestamp when the addon was created.

- `id` *string* &rarr;  The unique identifier of the cluster addon.

- `status` *string* &rarr;  The current deployment status of the addon (e.g., 'INSTALLING', 'UPDATING', 'DELETING').

- `updated_at` *string* &rarr;  The timestamp when the addon was last updated.



## Import

Cluster addons can be imported using the `id`, e.g.

```shell
terraform import vkcs_kubernetes_cluster_addon_v2.addon 3E4f0lX3N7HXNzyHl7YDTtX5pgY
```
4