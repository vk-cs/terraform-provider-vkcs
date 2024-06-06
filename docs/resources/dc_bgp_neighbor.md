---
subcategory: "Direct Connect"
layout: "vkcs"
page_title: "vkcs: vkcs_dc_bgp_neighbor"
description: |-
  Manages a direct connect BGP neighbor resource within VKCS.
---

# vkcs_dc_bgp_neighbor

Manages a direct connect BGP neighbor resource.

~> **Note:** This resource requires Sprut SDN to be enabled in your project.

**New since v0.5.0**.

## Example Usage
```terraform
resource "vkcs_dc_bgp_neighbor" "dc_bgp_neighbor" {
  name        = "tf-example"
  add_paths   = "on"
  description = "tf-example-description"
  dc_bgp_id   = vkcs_dc_bgp_instance.dc_bgp_instance.id
  remote_asn  = 1
  remote_ip   = "192.168.1.3"
}
```

## Argument Reference
- `dc_bgp_id` **required** *string* &rarr;  Direct Connect BGP ID to attach. Changing this creates a new resource

- `remote_asn` **required** *number* &rarr;  BGP Neighbor ASN. Changing this creates a new resource

- `remote_ip` **required** *string* &rarr;  BGP Neighbor IP address. Changing this creates a new resource

- `add_paths` optional *string* &rarr;  Activate BGP Add-Paths feature on peer. Default is off

- `bfd_enabled` optional *boolean* &rarr;  Control BGP session activity with BFD protocol. Default is false

- `description` optional *string* &rarr;  Description of the BGP neighbor

- `enabled` optional *boolean* &rarr;  Enable or disable item. Default is true

- `filter_in` optional *string* &rarr;  Input filter that pass incoming BGP prefixes (allow any)

- `filter_out` optional *string* &rarr;  Output filter that pass incoming BGP prefixes (allow any)

- `force_ibgp_next_hop_self` optional *boolean* &rarr;  Force set IP address of next-hop on BGP prefix to self even in iBGP. Default is false

- `name` optional *string* &rarr;  Name of the BGP neighbor

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `id` *string* &rarr;  ID of the resource

- `updated_at` *string* &rarr;  Update timestamp



## Import

Direct connect BGP neighbor can be imported using the `id`, e.g.
```shell
terraform import vkcs_dc_bgp_neighbor.mydcbgpneighbor 73096185-f200-4790-8095-962617b755f8
```
