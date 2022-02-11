---
layout: "vkcs"
page_title: "VKCS: compute_availability_zones"
description: |-
  Get a list of availability zones from OpenStack
---

# vkcs\_compute\_availability\_zones

Use this data source to get a list of availability zones from OpenStack

## Example Usage

```hcl
data "vkcs_compute_availability_zones" "zones" {}
```

## Argument Reference

* `region` - (Optional) The `region` to fetch availability zones from, defaults to the provider's `region`
* `state` - (Optional) The `state` of the availability zones to match, default ("available").


## Attributes Reference

`id` is set to hash of the returned zone list. In addition, the following attributes
are exported:

* `names` - The names of the availability zones, ordered alphanumerically, that match the queried `state`
