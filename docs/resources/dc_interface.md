---
subcategory: "Direct Connect"
layout: "vkcs"
page_title: "vkcs: vkcs_dc_interface"
description: |-
  Manages a direct connect interface resource within VKCS.
---

# vkcs_dc_interface

Manages a direct connect interface resource.

~> **Note:** This resource requires Sprut SDN to be enabled in your project.

**New since v0.5.0**.

## Example Usage
```terraform
# Connect networks to the router
resource "vkcs_dc_interface" "dc_interface" {
  name         = "tf-example"
  description  = "tf-example-description"
  dc_router_id = vkcs_dc_router.dc_router.id
  network_id   = vkcs_networking_network.app.id
  subnet_id    = vkcs_networking_subnet.app.id
}
```

## Connect dc router to Internet
```terraform
# Connect internet to the router
resource "vkcs_dc_interface" "dc_interface_internet" {
  name         = "interface-for-internet"
  dc_router_id = vkcs_dc_router.dc_router.id
  network_id   = data.vkcs_networking_network.internet_sprut.id
}
```

## Argument Reference
- `dc_router_id` **required** *string* &rarr;  Direct Connect Router ID to attach. Changing this creates a new resource

- `network_id` **required** *string* &rarr;  Network ID to attach. Changing this creates a new resource

- `bgp_announce_enabled` optional *boolean* &rarr;  Enable BGP announce of subnet attached to interface. Default is true

- `description` optional *string* &rarr;  Description of the interface

- `ip_address` optional *string* &rarr;  IP Address of the interface. Changing this creates a new resource

- `name` optional *string* &rarr;  Name of the interface

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`. Changing this creates a new interface.

- `subnet_id` optional *string* &rarr;  Subnet ID to attach. Changing this creates a new resource


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `id` *string* &rarr;  ID of the resource

- `ip_netmask` *number* &rarr;  IP Netmask

- `mac_address` *string* &rarr;  MAC Address of created interface

- `mtu` *number* &rarr;  MTU

- `port_id` *string* &rarr;  Port ID

- `sdn` *string* &rarr;  SDN where interface was created

- `updated_at` *string* &rarr;  Update timestamp



## Import

Direct connect interface can be imported using the `id`, e.g.
```shell
terraform import vkcs_dc_interface.mydcinterface 438d7479-d95f-4afc-b85e-eb8cd130a99f
```
