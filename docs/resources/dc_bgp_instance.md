---
subcategory: "Direct Connect"
layout: "vkcs"
page_title: "vkcs: vkcs_dc_bgp_instance"
description: |-
  Manages a direct connect BGP instance resource within VKCS.
---

# vkcs_dc_bgp_instance

Manages a direct connect BGP instance resource. **Note:** This resource requires Sprut SDN to be enabled in your project.

## Example Usage
```terraform
resource "vkcs_dc_bgp_instance" "dc_bgp_instance" {
    name = "tf-example"
    description = "tf-example-description"
    dc_router_id = vkcs_dc_router.dc_router.id
    bgp_router_id = "192.168.1.2"
    asn = 12345
    ecmp_enabled = true
    enabled = true
    graceful_restart = true
}
```

## Argument Reference
- `asn` **required** *number* &rarr;  BGP Autonomous System Number (integer representation supports only). Changing this creates a new resource

- `bgp_router_id` **required** *string* &rarr;  BGP Router ID (IP address that represent BGP router in BGP network). Changing this creates a new resource

- `dc_router_id` **required** *string* &rarr;  Direct Connect Router ID to attach. Changing this creates a new resource

- `description` optional *string* &rarr;  Description of the router

- `ecmp_enabled` optional *boolean* &rarr;  Enable BGP ECMP behaviour on router. Default is false

- `enabled` optional *boolean* &rarr;  Enable or disable item. Default is true

- `graceful_restart` optional *boolean* &rarr;  Enable BGP Graceful Restart feature. Default is false

- `long_lived_graceful_restart` optional *boolean* &rarr;  Enable BGP Long Lived Graceful Restart feature. Default is false

- `name` optional *string* &rarr;  Name of the router

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `id` *string* &rarr;  ID of the resource

- `updated_at` *string* &rarr;  Update timestamp



## Import

Direct connect BGP instance can be imported using the `name`, e.g.
```shell
terraform import vkcs_dc_bgp_instance.mydcbgpinstance e73496b2-e476-4536-9167-af24d18e1486
```
