---
layout: "vkcs"
page_title: "vkcs: vkcs_regions"
description: |-
  List available vkcs regions.
---

# vkcs_regions

`vkcs_regions` provides information about VKCS regions. Can be used to filter regions by parent region. To get details of each region the data source can be combined with the `vkcs_region` data source.

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
- `parent_region_id` **String** (*Optional*) ID of the parent region. Use empty value to list all the regions.


## Attributes Reference
- `parent_region_id` **String** See Argument Reference above.

- `id` **String** Random identifier of the data source.

- `names` <strong>Set of </strong>**String** Names of regions that meets the criteria.


