---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_clustertemplates"
description: |-
  List available vkcs kubernetes cluster templates.
---

# vkcs_kubernetes_clustertemplates

Use this data source to get a list of available VKCS Kubernetes Cluster Templates. To get details about each cluster template the data source can be combined with the `vkcs_kubernetes_clustertemplate` data source.

## Example Usage

Enabled VKCS Kubernetes Cluster Templates:
```terraform
data "vkcs_kubernetes_clustertemplates" "templates" {}

output "available_templates_by_name" {
  value = [
    for template in data.vkcs_kubernetes_clustertemplates.templates.cluster_templates :
    { name = template.name, version = template.version }
  ]
}
```
## Argument Reference
- `region` optional *string* &rarr;  The region to obtain the service client. If omitted, the `region` argument of the provider is used.<br>**New since v0.3.0**.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `cluster_templates`  *list* &rarr;  Available kubernetes cluster templates.
  - `cluster_template_uuid` *string* &rarr;  UUID of a cluster template.

  - `name` *string* &rarr;  Name of a cluster template.

  - `version` *string* &rarr;  Version of Kubernetes.


- `id` *string* &rarr;  Random identifier of the data source.


