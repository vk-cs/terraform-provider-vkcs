---
subcategory: "Kubernetes"
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_volume_types_v2"
description: |-
  Use this data source to retrieve information about available VKCS Kubernetes node group root volume types.
---

# vkcs_kubernetes_volume_types_v2

Use this data source to retrieve information about available VKCS Kubernetes node group root volume types.

## Example Usage
```terraform
data "vkcs_kubernetes_volume_types_v2" "available_volume_types" {}
```

## Argument Reference
- `region` optional *string* &rarr;  The region for which to retrieve volume types. Defaults to provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  A synthetic identifier set to "volume_types". This data source does not have a natural ID.

- `volume_types`  *set* &rarr;  A set of available volume types with their supported availability zones.
    - `name` *string* &rarr;  The name of the volume type (e.g., "ceph-hdd", "high-iops").

    - `zones` *set of* *string* &rarr;  A set of availability zones where this volume type is available (e.g., ["PA2", "MS1", "ME1"]).



