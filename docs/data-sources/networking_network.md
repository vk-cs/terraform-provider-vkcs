---
subcategory: "Network"
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
- `description` optional *string* &rarr;  Human-readable description of the network.

- `external` optional *boolean* &rarr;  The external routing facility of the network.

- `matching_subnet_cidr` optional *string* &rarr;  The CIDR of a subnet within the network.

- `name` optional *string* &rarr;  The name of the network.

- `network_id` optional *string* &rarr;  The ID of the network.

- `region` optional *string* &rarr;  The region in which to obtain the Network client. A Network client is needed to retrieve networks ids. If omitted, the `region` argument of the provider is used.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `status` optional *string* &rarr;  The status of the network.

- `tags` optional *set of* *string* &rarr;  The list of network tags to filter.

- `tenant_id` optional *string* &rarr;  The owner of the network.

- `vkcs_services_access` optional *boolean* &rarr;  Specifies whether VKCS services access is enabled.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `admin_state_up` *string* &rarr;  The administrative state of the network.

- `all_tags` *set of* *string* &rarr;  The set of string tags applied on the network.

- `id` *string* &rarr;  ID of the found network.

- `private_dns_domain` *string* &rarr;  Private dns domain name

- `shared` *string* &rarr;  Specifies whether the network resource can be accessed by any tenant or not.

- `subnets` *string* &rarr;  A list of subnet IDs belonging to the network.


