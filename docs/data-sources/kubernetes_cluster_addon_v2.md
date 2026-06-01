---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_cluster_addon_v2"
description: |-
  Use this data source to retrieve information about specific VKCS Kubernetes cluster addon.
---

# vkcs_kubernetes_cluster_addon_v2

Use this data source to retrieve information about specific VKCS Kubernetes cluster addon.

The cluster addon can be identified by either:

* `id`
* the combination of `base_addon_name` and `cluster_id`

These options are mutually exclusive.

## Example Usage

```terraform
data "vkcs_kubernetes_cluster_addon_v2" "ingress_nginx" {
  id = vkcs_kubernetes_cluster_addon_v2.ingress_nginx.id
}
```

## Argument Reference
- `base_addon_name` optional *string* &rarr;  The human-readable name of the Kubernetes addon (e.g., `ingress-nginx`, `argocd`).

- `cluster_id` optional *string* &rarr;  The ID of the Kubernetes cluster where the addon is installed.

- `id` optional *string* &rarr;  The unique identifier of the cluster addon.

- `region` optional *string* &rarr;  The region in which to get the cluster addon. If omitted, the provider's `region` is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `addon_id` *string* &rarr;  The unique identifier of the addon definition.

- `addon_name` *string* &rarr;  The human-readable name of the cluster addon.

- `addon_version_id` *string* &rarr;  The ID of the addon version to install.

- `created_at` *string* &rarr;  The timestamp when the addon was created.

- `namespace` *string* &rarr;  The Kubernetes namespace where the addon will be deployed.

- `status` *string* &rarr;  The current deployment status of the addon (e.g., 'INSTALLING', 'UPDATING', 'DELETING').

- `updated_at` *string* &rarr;  The timestamp when the addon was last updated.

- `values` *string* &rarr;  The Helm chart values for configuring the addon, typically provided as YAML string.


