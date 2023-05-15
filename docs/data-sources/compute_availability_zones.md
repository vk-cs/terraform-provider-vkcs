---
subcategory: "Virtual Machines"
layout: "vkcs"
page_title: "vkcs: vkcs_compute_availability_zones"
description: |-
  Get a list of availability zones from VKCS
---

# vkcs_compute_availability_zones

Use this data source to get a list of availability zones from VKCS

## Example Usage

```terraform
data "vkcs_compute_availability_zones" "zones" {}
```

## Argument Reference
- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`

- `state` optional *string* &rarr;  The `state` of the availability zones to match, default ("available").


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  Hash of the returned zone list.

- `names` *string* &rarr;  The names of the availability zones, ordered alphanumerically, that match the queried `state`


