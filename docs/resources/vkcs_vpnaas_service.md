---
layout: "vkcs"
page_title: "VKCS: vpnaas_service"
description: |-
	Manages a Neutron VPN service resource within OpenStack.
---

# vkcs\_vpnaas\_service

Manages a Neutron VPN service resource within OpenStack.

## Example Usage

```hcl
resource "vkcs_vpnaas_service" "service_1" {
	name           = "my_service"
	router_id      = "14a75700-fc03-4602-9294-26ee44f366b3"
	admin_state_up = "true"
}
```

## Argument Reference

The following arguments are supported:

* `admin_state_up` - (Optional) The administrative state of the resource. Can either be up(true) or down(false).
	Changing this updates the administrative state of the existing service.

* `description` - (Optional) The human-readable description for the service.
	Changing this updates the description of the existing service.

* `name` - (Optional) The name of the service. Changing this updates the name of
	the existing service.

* `region` - (Optional) The region in which to obtain the Networking client.
	A Networking client is needed to create a VPN service. If omitted, the
	`region` argument of the provider is used. Changing this creates a new
	service.

* `router_id` - (Required) The ID of the router. Changing this creates a new service.

* `subnet_id` - (Optional) SubnetID is the ID of the subnet. Default is null.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `router_id` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `subnet_id` - See Argument Reference above.
* `status` - Indicates whether IPsec VPN service is currently operational. Values are ACTIVE, DOWN, BUILD, ERROR, PENDING_CREATE, PENDING_UPDATE, or PENDING_DELETE.
* `external_v6_ip` - The read-only external (public) IPv6 address that is used for the VPN service.
* `external_v4_ip` - The read-only external (public) IPv4 address that is used for the VPN service.
* `description` - See Argument Reference above.

## Import

Services can be imported using the `id`, e.g.

```
$ terraform import vkcs_vpnaas_service.service_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
