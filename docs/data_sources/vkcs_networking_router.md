---
layout: "vkcs"
page_title: "VKCS: networking_router"
description: |-
  Get information on an OpenStack Floating IP.
---

# vkcs\_networking\_router

Use this data source to get the ID of an available OpenStack router.

## Example Usage

```hcl
data "vkcs_networking_router" "router" {
  name = "router_1"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve router ids. If omitted, the
  `region` argument of the provider is used.

* `router_id` - (Optional) The UUID of the router resource.

* `name` - (Optional) The name of the router.

* `description` - (Optional) Human-readable description of the router.

* `admin_state_up` - (Optional) Administrative up/down status for the router (must be "true" or "false" if provided).

* `status` - (Optional) The status of the router (ACTIVE/DOWN).

* `tags` - (Optional) The list of router tags to filter.

* `tenant_id` - (Optional) The owner of the router.

* `sdn` - (Optional) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

## Attributes Reference

`id` is set to the ID of the found router. In addition, the following attributes
are exported:

* `enable_snat` - The value that points out if the Source NAT is enabled on the router.

* `external_network_id` - The network UUID of an external gateway for the router.

* `all_tags` - The set of string tags applied on the router.

* `sdn` - See Argument Reference above.
