---
subcategory: "Direct Connect"
layout: "vkcs"
page_title: "vkcs: vkcs_dc_static_route"
description: |-
  Manages a direct connect static route resource within VKCS.
---

# vkcs_dc_static_route

Manages a direct connect BGP Static Announce resource.

~> **Note:** This resource requires Sprut SDN to be enabled in your project.

**New since v0.5.0**.

## Example Usage
```terraform
resource "vkcs_dc_static_route" "dc_static_route" {
  name         = "tf-example"
  description  = "tf-example-description"
  dc_router_id = vkcs_dc_router.dc_router.id
  network      = "192.168.1.0/24"
  gateway      = "192.168.1.3"
  metric       = 1
}
```

## Argument Reference
- `dc_router_id` **required** *string* &rarr;  Direct Connect Router ID to attach. Changing this creates a new resource

- `gateway` **required** *string* &rarr;  IP address of gateway. Changing this creates a new resource

- `network` **required** *string* &rarr;  Subnet in CIDR notation. Changing this creates a new resource

- `description` optional *string* &rarr;  Description of the static route

- `metric` optional *number* &rarr;  Metric to use for route. Default is 1

- `name` optional *string* &rarr;  Name of the static route

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `id` *string* &rarr;  ID of the resource

- `updated_at` *string* &rarr;  Update timestamp



## Import

Direct connect BGP instance can be imported using the `id`, e.g.
```shell
terraform import vkcs_dc_static_route.mydcstaticroute 2ee73dd1-d52a-4c3f-9041-c60900c154a4
```
