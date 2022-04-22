---
layout: "vkcs"
page_title: "kcs: kubernetes cluster templates"
description: |-
List available vkcs kubernetes cluster templates.
---

# VKCS Kubernetes Cluster Templates

`vkcs_kubernetes_cluster_templates` returns list of available VKCS Kubernetes Cluster Templates. 
To get details of each cluster template the data source can be combined with the `vkcs_kubernetes_clustertemplate` data source.

**New since version v0.5.0**

### Example Usage

Enabled VKCS Kubernetes Cluster Templates:

```hcl
data "vkcs_kubernetes_clustertemplates" "templates" {}
```

### Attributes Reference

* `id` - Random identifier of the data source.
* `cluster_templates` - A list of available kubernetes cluster templates.
  * `cluster_template_uuid` - The UUID of the cluster template.
  * `name` - The name of the cluster template.


