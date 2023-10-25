---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_addon"
description: |-
  Manages kubernetes cluster addons.
---

# vkcs_kubernetes_addon

Provides a kubernetes cluster addon resource. This can be used to create, modify and delete kubernetes cluster addons.

**New since v0.3.0**.

## Example Usage
```terraform
resource "vkcs_kubernetes_addon" "kube-prometheus-stack" {
  cluster_id           = vkcs_kubernetes_cluster.k8s-cluster.id
  addon_id             = data.vkcs_kubernetes_addon.kube-prometheus-stack.id
  namespace            = "monitoring"
  configuration_values = data.vkcs_kubernetes_addon.kube-prometheus-stack.configuration_values

  depends_on = [
    vkcs_kubernetes_node_group.default_ng
  ]
}
```

## Argument Reference
- `addon_id` **required** *string* &rarr;  The id of the addon. Changing this creates a new addon.

- `cluster_id` **required** *string* &rarr;  The ID of the kubernetes cluster. Changing this creates a new addon.

- `namespace` **required** *string* &rarr;  The namespace name where the addon will be installed.

- `configuration_values` optional *string* &rarr;  Configuration code for the addon. Changing this creates a new addon.

- `name` optional *string* &rarr;  The name of the application. Changing this creates a new addon.

- `region` optional *string* &rarr;  The region in which to obtain the Container Infra Addons client. If omitted, the `region` argument of the provider is used. Changing this creates a new addon.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource



## Import

Cluster addons can be imported using the `id` in the form `<cluster-id>/<cluster-addon-id>`, e.g.

```shell
terraform import vkcs_kubernetes_addon.addon 141a1f77-0e89-4b63-8d75-1b4ae496f862/a94c8ae2-0cac-4795-9253-d23ce2a70f86
```
