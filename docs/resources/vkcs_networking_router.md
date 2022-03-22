---
layout: "vkcs"
page_title: "VKCS: networking_router"
description: |-
  Manages a V2 router resource within OpenStack.
---

# vkcs\_networking\_router

Manages a V2 router resource within OpenStack.

## Example Usage

```hcl
resource "vkcs_networking_router" "router_1" {
  name                = "my_router"
  admin_state_up      = true
  external_network_id = "f67f0d72-0ddf-11e4-9d95-e1f29f417e2f"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 networking client.
  A networking client is needed to create a router. If omitted, the
  `region` argument of the provider is used. Changing this creates a new
  router.

* `name` - (Optional) A unique name for the router. Changing this
  updates the `name` of an existing router.

* `description` - (Optional) Human-readable description for the router.

* `admin_state_up` - (Optional) Administrative up/down status for the router
  (must be "true" or "false" if provided). Changing this updates the
  `admin_state_up` of an existing router.

* `external_network_id` - (Optional) The network UUID of an external gateway
  for the router. A router with an external gateway is required if any
  compute instances or load balancers will be using floating IPs. Changing
  this updates the external gateway of the router.

* `value_specs` - (Optional) Map of additional driver-specific options.

* `tags` - (Optional) A set of string tags for the router.

* `vendor_options` - (Optional) Map of additional vendor-specific options.
  Supported options are described below.

* `sdn` - (Optional) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

The `vendor_options` block supports:

* `set_router_gateway_after_create` - (Optional) Boolean to control whether
  the Router gateway is assigned during creation or updated after creation.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the router.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `external_network_id` - See Argument Reference above.
* `value_specs` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `all_tags` - The collection of tags assigned on the router, which have been
  explicitly and implicitly added.
* `sdn` - See Argument Reference above.

## Import

Routers can be imported using the `id`, e.g.

```
$ terraform import vkcs_networking_router.router_1 014395cd-89fc-4c9b-96b7-13d1ee79dad2
```
