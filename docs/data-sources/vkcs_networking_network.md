---
layout: "vkcs"
page_title: "vkcs: vkcs_networking_network"
description: |-
  Get information on an VKCS Network.
---

# vkcs_networking_network

Use this data source to get the ID of an available VKCS network.

## Example Usage

```terraform
data "vkcs_networking_network" "network" {
  name = "tf_test_network"
}
```

## Argument Reference
- `description` **String** (*Optional*) Human-readable description of the network.

- `external` **Boolean** (*Optional*) The external routing facility of the network.

- `matching_subnet_cidr` **String** (*Optional*) The CIDR of a subnet within the network.

- `name` **String** (*Optional*) The name of the network.

- `network_id` **String** (*Optional*) The ID of the network.

- `region` **String** (*Optional*) The region in which to obtain the Network client. A Network client is needed to retrieve networks ids. If omitted, the `region` argument of the provider is used.

- `sdn` **String** (*Optional*) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `status` **String** (*Optional*) The status of the network.

- `tags` <strong>Set of </strong>**String** (*Optional*) The list of network tags to filter.

- `tenant_id` **String** (*Optional*) The owner of the network.

- `vkcs_services_access` **Boolean** (*Optional*) Specifies whether VKCS services access is enabled.


## Attributes Reference
- `description` **String** See Argument Reference above.

- `external` **Boolean** See Argument Reference above.

- `matching_subnet_cidr` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `network_id` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `sdn` **String** See Argument Reference above.

- `status` **String** See Argument Reference above.

- `tags` <strong>Set of </strong>**String** See Argument Reference above.

- `tenant_id` **String** See Argument Reference above.

- `vkcs_services_access` **Boolean** See Argument Reference above.

- `admin_state_up` **String** The administrative state of the network.

- `all_tags` <strong>Set of </strong>**String** The set of string tags applied on the network.

- `id` **String** ID of the found network.

- `private_dns_domain` **String** Private dns domain name

- `shared` **String** Specifies whether the network resource can be accessed by any tenant or not.

- `subnets` **String** A list of subnet IDs belonging to the network.


