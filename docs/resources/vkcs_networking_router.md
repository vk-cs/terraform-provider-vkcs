---
layout: "vkcs"
page_title: "vkcs: vkcs_networking_router"
description: |-
  Manages a router resource within VKCS.
---

# vkcs_networking_router

Manages a router resource within VKCS.

## Example Usage
```terraform
resource "vkcs_networking_router" "router_1" {
  name                = "my_router"
  admin_state_up      = true
  external_network_id = "f67f0d72-0ddf-11e4-9d95-e1f29f417e2f"
}
```

## Argument Reference
- `admin_state_up` **Boolean** (*Optional*) Administrative up/down status for the router (must be "true" or "false" if provided). Changing this updates the `admin_state_up` of an existing router.

- `description` **String** (*Optional*) Human-readable description for the router.

- `external_network_id` **String** (*Optional*) The network UUID of an external gateway for the router. A router with an external gateway is required if any compute instances or load balancers will be using floating IPs. Changing this updates the external gateway of the router.

- `name` **String** (*Optional*) A unique name for the router. Changing this updates the `name` of an existing router.

- `region` **String** (*Optional*) The region in which to obtain the networking client. A networking client is needed to create a router. If omitted, the `region` argument of the provider is used. Changing this creates a new router.

- `sdn` **String** (*Optional*) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `tags` <strong>Set of </strong>**String** (*Optional*) A set of string tags for the router.

- `value_specs` <strong>Map of </strong>**String** (*Optional*) Map of additional driver-specific options.

- `vendor_options` (*Optional*) Map of additional vendor-specific options. Supported options are described below.
  - `set_router_gateway_after_create` **Boolean** (*Optional*) Boolean to control whether the Router gateway is assigned during creation or updated after creation.


## Attributes Reference
- `admin_state_up` **Boolean** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `external_network_id` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `sdn` **String** See Argument Reference above.

- `tags` <strong>Set of </strong>**String** See Argument Reference above.

- `value_specs` <strong>Map of </strong>**String** See Argument Reference above.

- `vendor_options`  See Argument Reference above.
  - `set_router_gateway_after_create` **Boolean** See Argument Reference above.

- `all_tags` <strong>Set of </strong>**String** The collection of tags assigned on the router, which have been explicitly and implicitly added.

- `id` **String** ID of the resource.



## Import

Routers can be imported using the `id`, e.g.

```shell
terraform import vkcs_networking_router.router_1 014395cd-89fc-4c9b-96b7-13d1ee79dad2
```
