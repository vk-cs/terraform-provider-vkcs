---
subcategory: "Regions,"
layout: "vkcs"
page_title: "vkcs: vkcs_regions"
description: |-
  List available vkcs regions.
---

# vkcs_regions

`vkcs_regions` provides information about VKCS regions. To get details of each region the data source can be combined with the `vkcs_region` data source.

## Example Usage

Enabled VKCS Regions:
```terraform
data "vkcs_regions" "current" {}
```

To see regions with the known Parent Region `parent_region_id` argument needs to be set.
```terraform
data "vkcs_regions" "current" {
  parent_region_id = "RegionOne"
}
```

## Argument Reference

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  Random identifier of the data source.

- `names` *set of* *string* &rarr;  Names of regions that meets the criteria.


