---
subcategory: "Direct Connect"
layout: "vkcs"
page_title: "vkcs: vkcs_dc_vrrp"
description: |-
  Manages a direct connect VRRP resource within VKCS.
---

# vkcs_dc_vrrp

Manages a direct connect VRRP resource.

~> **Note:** This resource requires Sprut SDN to be enabled in your project.

**New since v0.5.0**.

## Example Usage
```terraform
resource "vkcs_dc_vrrp" "dc_vrrp" {
    name = "tf-example"
    description = "tf-example-description"
    group_id = 100
    network_id = vkcs_networking_network.app.id
    subnet_id = vkcs_networking_subnet.app.id
    advert_interval = 1
}
```

## Argument Reference
- `group_id` **required** *number* &rarr;  VRRP Group ID

- `network_id` **required** *string* &rarr;  Network ID to attach. Changing this creates a new resource

- `advert_interval` optional *number* &rarr;  VRRP Advertise interval. Default is 1

- `description` optional *string* &rarr;  Description of the VRRP

- `enabled` optional *boolean* &rarr;  Enable or disable item. Default is true

- `name` optional *string* &rarr;  Name of the VRRP

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.

- `subnet_id` optional *string* &rarr;  Subnet ID to attach. Changing this creates a new resource


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `id` *string* &rarr;  ID of the resource

- `sdn` *string* &rarr;  SDN of created VRRP

- `updated_at` *string* &rarr;  Update timestamp



## Import

Direct connect vrrp can be imported using the `name`, e.g.
```shell
terraform import vkcs_dc_vrrp.mydcvrrp f6149e79-b441-4327-90fc-7653acbc204c
```
