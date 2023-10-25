---
subcategory: "Direct Connect"
layout: "vkcs"
page_title: "vkcs: vkcs_dc_bgp_static_announce"
description: |-
  Manages a direct connect BGP static announce resource within VKCS.
---

# vkcs_dc_bgp_static_announce

Manages a direct connect BGP Static Announce resource.

~> **Note:** This resource requires Sprut SDN to be enabled in your project.

**New since v0.5.0**.

## Example Usage
```terraform
resource "vkcs_dc_bgp_static_announce" "dc_bgp_static_announce" {
    name = "tf-example"
    description = "tf-example-description"
    dc_bgp_id = vkcs_dc_bgp_instance.dc_bgp_instance.id
    network = "192.168.1.0/24"
    gateway = "192.168.1.3"
}
```

## Argument Reference
- `dc_bgp_id` **required** *string* &rarr;  Direct Connect BGP ID to attach. Changing this creates a new resource

- `gateway` **required** *string* &rarr;  IP address of gateway. Changing this creates a new resource

- `network` **required** *string* &rarr;  Subnet in CIDR notation. Changing this creates a new resource

- `description` optional *string* &rarr;  Description of the BGP neighbor

- `enabled` optional *boolean* &rarr;  Enable or disable item. Default is true

- `name` optional *string* &rarr;  Name of the BGP neighbor

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `id` *string* &rarr;  ID of the resource

- `updated_at` *string* &rarr;  Update timestamp



## Import

Direct connect BGP instance can be imported using the `name`, e.g.
```shell
terraform import vkcs_dc_bgp_static_announce.mydcbgpstaticannounce 8a1d9812-305b-468f-8ae5-833e181b01a8
```
