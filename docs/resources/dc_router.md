---
subcategory: "Direct Connect"
layout: "vkcs"
page_title: "vkcs: vkcs_dc_router"
description: |-
  Manages a direct connect router resource within VKCS.
---

# vkcs_dc_router

Manages a direct connect router resource. <br> ~> **Note:** This resource requires Sprut SDN to be enabled in your project. **New since v0.5.0**.

## Example Usage
```terraform
resource "vkcs_dc_router" "dc_router" {
  availability_zone = "GZ1"
  flavor = "standard"
  name = "tf-example"
  description = "tf-example-description"
}
```

## Argument Reference
- `availability_zone` optional *string* &rarr;  The availability zone in which to create the router. Changing this creates a new router

- `description` optional *string* &rarr;  Description of the router

- `flavor` optional *string* &rarr;  Flavor of the router. Possible values can be obtained with vkcs_dc_api_options data source. Changing this creates a new router. <br>**Note:** Not to be confused with compute service flavors.

- `name` optional *string* &rarr;  Name of the router

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `id` *string* &rarr;  ID of the resource

- `updated_at` *string* &rarr;  Update timestamp



## Import

Direct connect router can be imported using the `name`, e.g.
```shell
terraform import vkcs_dc_router.mydcrouter b50b32fc-16e2-4cb0-acdb-638c865c4242
```
