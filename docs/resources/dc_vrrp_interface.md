---
subcategory: "Direct Connect"
layout: "vkcs"
page_title: "vkcs: vkcs_dc_vrrp_interface"
description: |-
  Manages a direct connect VRRP interface resource within VKCS.
---

# vkcs_dc_vrrp_interface

Manages a direct connect VRRP interface resource.

## Example Usage
```terraform
resource "vkcs_dc_vrrp_interface" "dc_vrrp_interface" {
    name = "tf-example"
    description = "tf-example-description"
    dc_vrrp_id = vkcs_dc_vrrp.dc_vrrp.id
    dc_interface_id = vkcs_dc_interface.dc_interface.id
    priority = 100
    preempt = true
    master = true
}
```

## Argument Reference
- `dc_interface_id` **required** *string* &rarr;  DC Interface ID to attach. Changing this creates a new resource

- `dc_vrrp_id` **required** *string* &rarr;  VRRP ID to attach. Changing this creates a new resource

- `description` optional *string* &rarr;  Description of the VRRP

- `master` optional *boolean* &rarr;  Start VRRP instance on interface as VRRP Master. Default is false

- `name` optional *string* &rarr;  Name of the VRRP

- `preempt` optional *boolean* &rarr;  VRRP interface preempt behavior. Default is true

- `priority` optional *number* &rarr;  VRRP interface priority. Default is 100

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `id` *string* &rarr;  ID of the resource

- `updated_at` *string* &rarr;  Update timestamp



## Import

Direct connect vrrp interface can be imported using the `name`, e.g.
```shell
terraform import vkcs_dc_vrrp_interface.mydcvrrpinterface 3f071a6d-3d21-435c-83f7-11b276f318f0
```
