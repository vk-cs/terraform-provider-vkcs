---
subcategory: "Network"
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
  name = "router-tf-example"
  tags = ["tf-example"]
  # This is unnecessary in real life.
  # This is required here to let the example work with router resource example. 
  depends_on = [vkcs_networking_router.router]
}
```

## Argument Reference
- `admin_state_up` optional *boolean* &rarr;  Administrative up/down status for the router (must be "true" or "false" if provided).

- `description` optional *string* &rarr;  Human-readable description of the router.

- `enable_snat` optional *boolean* &rarr;  The value that points out if the Source NAT is enabled on the router.

- `id` optional *string* &rarr;  The UUID of the router resource.

- `name` optional *string* &rarr;  The name of the router.

- `region` optional *string* &rarr;  The region in which to obtain the Network client. A Network client is needed to retrieve router ids. If omitted, the `region` argument of the provider is used.

- `router_id` optional deprecated *string* &rarr;  The UUID of the router resource. **Deprecated** This argument is deprecated, please, use the `id` attribute instead.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is project's default SDN.

- `status` optional *string* &rarr;  The status of the router (ACTIVE/DOWN).

- `tags` optional *set of* *string* &rarr;  The list of router tags to filter.

- `tenant_id` optional *string* &rarr;  The owner of the router.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `all_tags` *set of* *string* &rarr;  The set of string tags applied on the router.

- `external_fixed_ips` *object* &rarr;  List of external gateways of the router.<br>**New since v0.7.4**.

- `external_network_id` *string* &rarr;  The network UUID of an external gateway for the router.


