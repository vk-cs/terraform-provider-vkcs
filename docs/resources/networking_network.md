---
subcategory: "Network"
layout: "vkcs"
page_title: "vkcs: vkcs_networking_network"
description: |-
  Manages a network resource within VKCS.
---

# vkcs_networking_network

Manages a network resource within VKCS.

## Example Usage
```terraform
resource "vkcs_networking_network" "app" {
  name        = "app-tf-example"
  description = "Application network"
}
```
## Argument Reference
- `admin_state_up` optional *boolean* &rarr;  The administrative state of the network. Acceptable values are "true" and "false". Changing this value updates the state of the existing network.

- `description` optional *string* &rarr;  Human-readable description of the network. Changing this updates the name of the existing network.

- `name` optional *string* &rarr;  The name of the network. Changing this updates the name of the existing network.

- `port_security_enabled` optional *boolean* &rarr;  Whether to explicitly enable or disable port security on the network. Port Security is usually enabled by default, so omitting this argument will usually result in a value of "true". Setting this explicitly to `false` will disable port security. Valid values are `true` and `false`.

- `private_dns_domain` optional *string* &rarr;  Private dns domain name

- `region` optional *string* &rarr;  The region in which to obtain the Networking client. A Networking client is needed to create a network. If omitted, the `region` argument of the provider is used. Changing this creates a new network.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is project's default SDN.

- `tags` optional *set of* *string* &rarr;  A set of string tags for the network.

- `value_specs` optional *map of* *string* &rarr;  Map of additional options.

- `vkcs_services_access` optional *boolean* &rarr;  Whether VKCS services access is enabled. This feature should be enabled globally for your project. Access can be enabled for new or existing networks, but cannot be disabled for existing networks. Valid values are `true` and `false`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `all_tags` *set of* *string* &rarr;  The collection of tags assigned on the network, which have been explicitly and implicitly added.

- `id` *string* &rarr;  ID of the resource.



## Import

Networks can be imported using the `id`, e.g.

```shell
terraform import vkcs_networking_network.network_1 d90ce693-5ccf-4136-a0ed-152ce412b6b9
```
