---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_addons_v2"
description: |-
  Use this data source to retrieve information about available VKCS Kubernetes addons.
---

# vkcs_kubernetes_addons_v2

Use this data source to retrieve information about available VKCS Kubernetes addons.

## Example Usage

```terraform
data "vkcs_kubernetes_addons_v2" "kubernetes_addons" {}
```

## Argument Reference
- `region` optional *string* &rarr;  The region for which to retrieve list of addons. Defaults to provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `addons`  *set* &rarr;  List of available addons.
    - `id` *string* &rarr;  Unique identifier of the addon.

    - `name` *string* &rarr;  Human-readable name of the addon.

    - `versions`  *list* &rarr;  List of available versions for the addon.
        - `id` *string* &rarr;  Unique identifier of the addon version.

        - `version` *string* &rarr;  Version string of the addon.



- `id` *string* &rarr;  A synthetic identifier for the data source.


