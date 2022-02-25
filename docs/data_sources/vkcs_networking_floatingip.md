---
layout: "vkcs"
page_title: "VKCS: networking_floatingip"
description: |-
  Get information on an OpenStack Floating IP.
---

# vkcs\_networking\_floatingip

Use this data source to get the ID of an available OpenStack floating IP.

## Example Usage

```hcl
data "vkcs_networking_floatingip" "floatingip_1" {
  address = "192.168.0.4"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Neutron client.
  A Neutron client is needed to retrieve floating IP ids. If omitted, the
  `region` argument of the provider is used.

* `description` - (Optional) Human-readable description of the floating IP.

* `address` - (Optional) The IP address of the floating IP.

* `pool` - (Optional) The name of the pool from which the floating IP belongs to.

* `port_id` - (Optional) The ID of the port the floating IP is attached.

* `status` - status of the floating IP (ACTIVE/DOWN).

* `fixed_ip` - (Optional) The specific IP address of the internal port which should be associated with the floating IP.

* `tags` - (Optional) The list of floating IP tags to filter.

* `tenant_id` - (Optional) The owner of the floating IP.

## Attributes Reference

`id` is set to the ID of the found floating IP. In addition, the following attributes
are exported:

* `all_tags` - A set of string tags applied on the floating IP.
