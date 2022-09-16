---
layout: "vkcs"
page_title: "vkcs: vkcs_networking_router"
description: |-
  Get information on a VKCS router.
---

# vkcs_networking_router

Use this data source to get the ID of an available VKCS router.

## Example Usage

```terraform
data "vkcs_networking_router" "router" {
  name = "router_1"
}
```

## Argument Reference
- `admin_state_up` **Boolean** (*Optional*) Administrative up/down status for the router (must be "true" or "false" if provided).

- `description` **String** (*Optional*) Human-readable description of the router.

- `enable_snat` **Boolean** (*Optional*) The value that points out if the Source NAT is enabled on the router.

- `name` **String** (*Optional*) The name of the router.

- `region` **String** (*Optional*) The region in which to obtain the Network client. A Network client is needed to retrieve router ids. If omitted, the `region` argument of the provider is used.

- `router_id` **String** (*Optional*) The UUID of the router resource.

- `sdn` **String** (*Optional*) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `status` **String** (*Optional*) The status of the router (ACTIVE/DOWN).

- `tags` <strong>Set of </strong>**String** (*Optional*) The list of router tags to filter.

- `tenant_id` **String** (*Optional*) The owner of the router.


## Attributes Reference
- `admin_state_up` **Boolean** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `enable_snat` **Boolean** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `router_id` **String** See Argument Reference above.

- `sdn` **String** See Argument Reference above.

- `status` **String** See Argument Reference above.

- `tags` <strong>Set of </strong>**String** See Argument Reference above.

- `tenant_id` **String** See Argument Reference above.

- `all_tags` <strong>Set of </strong>**String** The set of string tags applied on the router.

- `external_network_id` **String** The network UUID of an external gateway for the router.

- `id` **String** ID of the found router.


