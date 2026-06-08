---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_addon_v2"
description: |-
  Use this data source to retrieve information about specific VKCS Kubernetes addon.
---

# vkcs_kubernetes_addon_v2

Use this data source to retrieve information about specific VKCS Kubernetes addon.

The addon can be identified by either:

* `id`
* the combination of `name` and `version`

These options are mutually exclusive.

## Example Usage

```terraform
data "vkcs_kubernetes_addon_v2" "ingress_nginx" {
  name    = "ingress-nginx"
  version = "4.12.1"
}
```

## Argument Reference
- `id` optional *string* &rarr;  Unique identifier of the addon version.

- `name` optional *string* &rarr;  Name of the addon.

- `region` optional *string* &rarr;  The region for which to retrieve specific addon. Defaults to provider's `region`.

- `version` optional *string* &rarr;  Version of the addon.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `addon_id` *string* &rarr;  Identifier of the addon this version belongs to.

- `supported_kube_versions` *set of* *string* &rarr;  List of supported Kubernetes versions.

- `values_template` *string* &rarr;  Base64-encoded values.yaml template for the addon.


