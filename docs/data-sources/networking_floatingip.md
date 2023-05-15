---
subcategory: "Network"
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
- `address` optional *string* &rarr;  The IP address of the floating IP.

- `description` optional *string* &rarr;  Human-readable description of the floating IP.

- `fixed_ip` optional *string* &rarr;  The specific IP address of the internal port which should be associated with the floating IP.

- `pool` optional *string* &rarr;  The name of the pool from which the floating IP belongs to.

- `port_id` optional *string* &rarr;  The ID of the port the floating IP is attached.

- `region` optional *string* &rarr;  The region in which to obtain the Network client. A Network client is needed to retrieve floating IP ids. If omitted, the `region` argument of the provider is used.

- `sdn` optional *string* &rarr;  SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

- `status` optional *string* &rarr;  Status of the floating IP (ACTIVE/DOWN).

- `tenant_id` optional *string* &rarr;  The owner of the floating IP.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the found floating IP.


