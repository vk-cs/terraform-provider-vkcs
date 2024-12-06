---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_addon"
description: |-
  Get information on a kubernetes cluster addon.
---

# vkcs_kubernetes_addon

Provides a kubernetes cluster addon datasource. This can be used to get information about an VKCS cluster addon.

**New since v0.3.0**.

## Example Usage

```terraform
data "vkcs_kubernetes_addon" "kube_prometheus_stack" {
  cluster_id = vkcs_kubernetes_cluster.k8s_cluster.id
  name       = "kube-prometheus-stack"
  version    = "54.2.2"
}
```

## Argument Reference
- `cluster_id` **required** *string* &rarr;  The ID of the kubernetes cluster.

- `name` **required** *string* &rarr;  An addon name to filter by.

- `version` **required** *string* &rarr;  An addon version to filter by.

- `region` optional *string* &rarr;  The region in which to obtain the service client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `configuration_values` *string* &rarr;  Configuration code for the addon. If the addon was installed in the cluster, this value is the user-provided configuration code, otherwise it is a template for this cluster.

- `id` *string* &rarr;  ID of the resource.

- `installed` *boolean* &rarr;  Whether the addon was installed in the cluster.


