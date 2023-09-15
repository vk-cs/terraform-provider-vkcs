---
subcategory: "Direct Connect"
layout: "vkcs"
page_title: "vkcs: vkcs_dc_api_options"
description: |-
  Get information on an VKCS Direct Connect API Options.
---

# vkcs_dc_api_options

Use this data source to get direct connect api options. **Note:** This resource requires Sprut SDN to be enabled in your project.

## Example Usage

```terraform
data "vkcs_dc_api_options" "dc_api_options" {}
```

## Argument Reference
- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `availability_zones` *string* &rarr;  List of avalability zone options

- `flavors` *string* &rarr;  List of flavor options for vkcs_dc_router resource

- `id` *string* &rarr;  ID of the data source


