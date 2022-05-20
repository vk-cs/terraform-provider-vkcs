---
layout: "vkcs"
page_title: "vkcs: regions"
description: |-
List available vkcs regions.
---

`vkcs_regions` provides information about VKCS regions. Can be used to filter regions by parent region. To get details of each region the data source can be combined with the `vkcs_region` data source.

### Example Usage

Enabled VKCS Regions:

```hcl
data "vkcs_regions" "current" {}
```

To see regions with the known Parent Region `parent_region_id` argument needs to be set.

```hcl
data "vkcs_regions" "current" {
  parent_region_id = "RegionOne"
}
```

### Argument Reference

The following arguments are supported:

* `parent_region_id` - (Optional) ID of the parent region. Use empty value to list all the regions.

### Attributes Reference

* `id` - Random identifier of the data source.
* `names` - Names of regions that meets the criteria.
