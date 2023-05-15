---
subcategory: "Network"
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
- `admin_state_up` optional *boolean* &rarr;  Administrative up/down status for the router (must be "true" or "false" if provided). Changing this updates the `admin_state_up` of an existing router.

- `description` optional *string* &rarr;  Human-readable description for the router.

- `external_network_id` optional *string* &rarr;  The network UUID of an external gateway for the router. A router with an external gateway is required if any compute instances or load balancers will be using floating IPs. Changing this updates the external gateway of the router.

- `name` optional *string* &rarr;  A unique name for the router. Changing this updates the `name` of an existing router.

- `region` optional *string* &rarr;  The region in which to obtain the networking client. A networking client is needed to create a router. If omitted, the `region` argument of the provider is used. Changing this creates a new router.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `tags` optional *set of* *string* &rarr;  A set of string tags for the router.

- `value_specs` optional *map of* *string* &rarr;  Map of additional driver-specific options.

- `vendor_options` optional &rarr;  Map of additional vendor-specific options. Supported options are described below.
  - `set_router_gateway_after_create` optional *boolean* &rarr;  Boolean to control whether the Router gateway is assigned during creation or updated after creation.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `all_tags` *set of* *string* &rarr;  The collection of tags assigned on the router, which have been explicitly and implicitly added.

- `id` *string* &rarr;  ID of the resource.



## Import

Routers can be imported using the `id`, e.g.

```shell
terraform import vkcs_networking_router.router_1 014395cd-89fc-4c9b-96b7-13d1ee79dad2
```
