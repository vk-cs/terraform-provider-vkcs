---
layout: "vkcs"
page_title: "vkcs: vkcs_vpnaas_endpoint_group"
description: |-
  Manages an Endpoint Group resource within VKCS.
---

# vkcs_vpnaas_endpoint_group

Manages an Endpoint Group resource within VKCS.

## Example Usage
```terraform
resource "vkcs_vpnaas_endpoint_group" "group_1" {
	name = "Group 1"
	type = "cidr"
	endpoints = [
		"10.2.0.0/24",
		"10.3.0.0/24",
	]
}
```
## Argument Reference
- `description` **String** (*Optional*) The human-readable description for the group. Changing this updates the description of the existing group.

- `endpoints` <strong>Set of </strong>**String** (*Optional*) List of endpoints of the same type, for the endpoint group. The values will depend on the type. Changing this creates a new group.

- `name` **String** (*Optional*) The name of the group. Changing this updates the name of the existing group.

- `region` **String** (*Optional*) The region in which to obtain the Networking client. A Networking client is needed to create an endpoint group. If omitted, the `region` argument of the provider is used. Changing this creates a new group.

- `type` **String** (*Optional*) The type of the endpoints in the group. A valid value is subnet, cidr, network, router, or vlan. Changing this creates a new group.


## Attributes Reference
- `description` **String** See Argument Reference above.

- `endpoints` <strong>Set of </strong>**String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `type` **String** See Argument Reference above.

- `id` **String** ID of the resource.



## Import

Groups can be imported using the `id`, e.g.

```shell
terraform import vkcs_vpnaas_endpoint_group.group_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
