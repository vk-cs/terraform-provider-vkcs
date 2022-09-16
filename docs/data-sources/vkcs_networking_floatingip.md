---
layout: "vkcs"
page_title: "vkcs: vkcs_networking_floatingip"
description: |-
  Get information on an VKCS Floating IP.
---

# vkcs_networking_floatingip

Use this data source to get the ID of an available VKCS floating IP.

## Example Usage

```terraform
data "vkcs_networking_floatingip" "floatingip_1" {
  address = "192.168.0.4"
}
```

## Argument Reference
- `address` **String** (*Optional*) The IP address of the floating IP.

- `description` **String** (*Optional*) Human-readable description of the floating IP.

- `fixed_ip` **String** (*Optional*) The specific IP address of the internal port which should be associated with the floating IP.

- `pool` **String** (*Optional*) The name of the pool from which the floating IP belongs to.

- `port_id` **String** (*Optional*) The ID of the port the floating IP is attached.

- `region` **String** (*Optional*) The region in which to obtain the Network client. A Network client is needed to retrieve floating IP ids. If omitted, the `region` argument of the provider is used.

- `sdn` **String** (*Optional*) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `status` **String** (*Optional*) Status of the floating IP (ACTIVE/DOWN).

- `tenant_id` **String** (*Optional*) The owner of the floating IP.


## Attributes Reference
- `address` **String** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `fixed_ip` **String** See Argument Reference above.

- `pool` **String** See Argument Reference above.

- `port_id` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `sdn` **String** See Argument Reference above.

- `status` **String** See Argument Reference above.

- `tenant_id` **String** See Argument Reference above.

- `id` **String** ID of the found floating IP.


