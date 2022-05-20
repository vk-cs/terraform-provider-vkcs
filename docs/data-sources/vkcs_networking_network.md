---
layout: "vkcs"
page_title: "vkcs: networking_network"
description: |-
  Get information on an VKCS Network.
---

# vkcs\_networking\_network

Use this data source to get the ID of an available VKCS network.

## Example Usage

```hcl
data "vkcs_networking_network" "network" {
  name = "tf_test_network"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the Network client.
  A Network client is needed to retrieve networks ids. If omitted, the
  `region` argument of the provider is used.

* `network_id` - (Optional) The ID of the network.

* `name` - (Optional) The name of the network.

* `description` - (Optional) Human-readable description of the network.

* `status` - (Optional) The status of the network.

* `external` - (Optional) The external routing facility of the network.

* `matching_subnet_cidr` - (Optional) The CIDR of a subnet within the network.

* `tenant_id` - (Optional) The owner of the network.

* `tags` - (Optional) The list of network tags to filter.

* `sdn` - (Optional) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

## Attributes Reference

`id` is set to the ID of the found network. In addition, the following attributes
are exported:

* `admin_state_up` - The administrative state of the network.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `region` - See Argument Reference above.
* `external` - See Argument Reference above.
* `shared` - Specifies whether the network resource can be accessed by any
   tenant or not.
* `mtu` - See Argument Reference above.
* `subnets` - A list of subnet IDs belonging to the network.
* `all_tags` - The set of string tags applied on the network.
* `private_dns_domain` - See Argument Reference above.
* `sdn` - See Argument Reference above.
* `vkcs_services_access` - Specifies whether VKCS services access is enabled.
