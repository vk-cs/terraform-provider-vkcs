---
subcategory: "Direct Connect"
layout: "vkcs"
page_title: "vkcs: vkcs_dc_vrrp_address"
description: |-
  Manages a direct connect VRRP address resource within VKCS.
---

# vkcs_dc_vrrp_address

Manages a direct connect VRRP address resource.

~> **Note:** This resource requires Sprut SDN to be enabled in your project.

**New since v0.5.0**.

## Example Usage
```terraform
resource "vkcs_dc_vrrp_address" "dc_vrrp_address" {
    name = "tf-example"
    description = "tf-example-description"
    dc_vrrp_id = vkcs_dc_vrrp.dc_vrrp.id
    ip_address = "192.168.199.42"
}
```

## Argument Reference
- `dc_vrrp_id` **required** *string* &rarr;  VRRP ID to attach. Changing this creates a new resource

- `description` optional *string* &rarr;  Description of the VRRP

- `ip_address` optional *string* &rarr;  IP address to assign. Changing this creates a new resource

- `name` optional *string* &rarr;  Name of the VRRP

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `id` *string* &rarr;  ID of the resource

- `port_id` *string* &rarr;  Port ID used to assign IP address

- `updated_at` *string* &rarr;  Update timestamp



## Import

Direct connect vrrp address can be imported using the `name`, e.g.
```shell
terraform import vkcs_dc_vrrp_address.mydcvrrpaddress aa00d2a9-db9c-4976-898b-fcabb9f49505
```
