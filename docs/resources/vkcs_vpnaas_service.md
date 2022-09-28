---
layout: "vkcs"
page_title: "vkcs: vkcs_vpnaas_service"
description: |-
  Manages a VPN service resource within VKCS.
---

# vkcs_vpnaas_service

Manages a VPN service resource within VKCS.

## Example Usage
```terraform
resource "vkcs_vpnaas_service" "service_1" {
	name           = "my_service"
	router_id      = "14a75700-fc03-4602-9294-26ee44f366b3"
	admin_state_up = "true"
}
```
## Argument Reference
- `router_id` **String** (***Required***) The ID of the router. Changing this creates a new service.

- `admin_state_up` **Boolean** (*Optional*) The administrative state of the resource. Can either be up(true) or down(false). Changing this updates the administrative state of the existing service.

- `description` **String** (*Optional*) The human-readable description for the service. Changing this updates the description of the existing service.

- `name` **String** (*Optional*) The name of the service. Changing this updates the name of the existing service.

- `region` **String** (*Optional*) The region in which to obtain the Networking client. A Networking client is needed to create a VPN service. If omitted, the `region` argument of the provider is used. Changing this creates a new service.

- `subnet_id` **String** (*Optional*) SubnetID is the ID of the subnet. Default is null.


## Attributes Reference
- `router_id` **String** See Argument Reference above.

- `admin_state_up` **Boolean** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `subnet_id` **String** See Argument Reference above.

- `external_v4_ip` **String** The read-only external (public) IPv4 address that is used for the VPN service.

- `external_v6_ip` **String** The read-only external (public) IPv6 address that is used for the VPN service.

- `id` **String** ID of the resource.

- `status` **String** Indicates whether IPsec VPN service is currently operational. Values are ACTIVE, DOWN, BUILD, ERROR, PENDING_CREATE, PENDING_UPDATE, or PENDING_DELETE.



## Import

Services can be imported using the `id`, e.g.

```shell
terraform import vkcs_vpnaas_service.service_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
