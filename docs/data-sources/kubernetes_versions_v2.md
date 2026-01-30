---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_versions_v2"
description: |-
  Returns a list of Kubernetes versions available for provisioning next-generation clusters.
---

# vkcs_kubernetes_versions_v2



## Example Usage
```terraform
data "vkcs_kubernetes_versions_v2" "available_versions" {}

output "available_kubernetes_versions" {
  description = "A set of Kubernetes versions that can be used to deploy a new cluster."
  value       = data.vkcs_kubernetes_versions_v2.available_versions.k8s_versions
}
```

## Argument Reference
- `region` optional *string* &rarr;  The region for which to retrieve available versions. Defaults to provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  A synthetic identifier set to "kubernetes_versions". This data source does not have a natural ID.

- `k8s_versions` *set of* *string* &rarr;  A set of available Kubernetes version strings (e.g., ["v1.33.3", "v1.34.2"]).


