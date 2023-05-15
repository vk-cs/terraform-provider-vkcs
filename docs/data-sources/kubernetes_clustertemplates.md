---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_clustertemplates"
description: |-
  List available vkcs kubernetes cluster templates.
---

# vkcs_kubernetes_clustertemplates

`vkcs_kubernetes_cluster_templates` returns list of available VKCS Kubernetes Cluster Templates. To get details of each cluster template the data source can be combined with the `vkcs_kubernetes_clustertemplate` data source.

## Example Usage

Enabled VKCS Kubernetes Cluster Templates:
```terraform
data "vkcs_kubernetes_clustertemplates" "templates" {}
```
## Argument Reference

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `cluster_templates` *object* &rarr;  A list of available kubernetes cluster templates.
  - `cluster_template_uuid` **String** The UUID of the cluster template.

  - `name` **String** The name of the cluster template.

  - `version` **String** The version of the cluster template.

- `id` *string* &rarr;  Random identifier of the data source.


