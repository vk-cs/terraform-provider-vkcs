---
layout: "vkcs"
page_title: "vkcs: networking_floatingip"
description: |-
  Get information on an VKCS Floating IP.
---

# vkcs\_networking\_floatingip

Use this data source to get the ID of an available VKCS floating IP.

## Example Usage

```hcl
data "vkcs_networking_floatingip" "floatingip_1" {
  address = "192.168.0.4"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the Network client.
  A Network client is needed to retrieve floating IP ids. If omitted, the
  `region` argument of the provider is used.

* `description` - (Optional) Human-readable description of the floating IP.

* `address` - (Optional) The IP address of the floating IP.

* `pool` - (Optional) The name of the pool from which the floating IP belongs to.

* `port_id` - (Optional) The ID of the port the floating IP is attached.

* `status` - status of the floating IP (ACTIVE/DOWN).

* `fixed_ip` - (Optional) The specific IP address of the internal port which should be associated with the floating IP.

* `tenant_id` - (Optional) The owner of the floating IP.

* `sdn` - (Optional) SDN to use for this resource. Must be one of following: "neutron", "sprut". Default value is "neutron".

## Attributes Reference

`id` is set to the ID of the found floating IP. In addition, the following attributes
are exported:

* `sdn` - See Argument Reference above.
