---
subcategory: "Regions,"
layout: "vkcs"
page_title: "vkcs: vkcs_region"
description: |-
  Get information about region.
---

# vkcs_region

`vkcs_region` provides details about a specific VKCS region. As well as validating a given region name this resource can be used to discover the name of the region configured within the provider.

## Example Usage

The following example shows how the resource might be used to obtain the name of the VKCS region configured on the provider.

```terraform
data "vkcs_region" "current" {}
```

## Argument Reference
- `id` optional *string* &rarr;  ID of the region to learn or use. Use empty value to learn current region on the provider.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `description` *string* &rarr;  Description of the region.

- `parent_region` *string* &rarr;  Parent of the region.


