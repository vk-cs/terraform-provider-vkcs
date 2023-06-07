---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_addons"
description: |-
  Get information on kubernetes cluster addons.
---

# vkcs_kubernetes_addons

Provides a kubernetes cluster addons datasource. This can be used to get information about an VKCS cluster addons.

## Example Usage

```terraform
data "vkcs_kubernetes_addons" "cluster-addons" {
  cluster_id = vkcs_kubernetes_cluster.k8s-cluster.id
}
```

## Argument Reference
- `cluster_id` **required** *string* &rarr;  The ID of the kubernetes cluster.

- `region` optional *string* &rarr;  The region in which to obtain the service client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `addons`  *list*
  - `id` *string* &rarr;  ID of an addon.

  - `installed` *boolean* &rarr;  Whether an addon was installed in the cluster.

  - `name` *string* &rarr;  Name of an addon.

  - `version` *string* &rarr;  Version of an addon.


- `id` *string* &rarr;  ID of the resource.


