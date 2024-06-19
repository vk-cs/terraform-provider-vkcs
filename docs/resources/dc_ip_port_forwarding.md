---
subcategory: "Direct Connect"
layout: "vkcs"
page_title: "vkcs: vkcs_dc_ip_port_forwarding"
description: |-
  Manages a direct connect IP port forwarding resource within VKCS.
---

# vkcs_dc_ip_port_forwarding

Manages a direct connect ip port forwarding resource.

~> **Note:** This resource requires Sprut SDN to be enabled in your project.

**New since v0.8.0**.

## Example Usage
```terraform
resource "vkcs_dc_ip_port_forwarding" "dc-ip-port-forwarding" {
  dc_interface_id = vkcs_dc_interface.dc_interface.id
  name            = "tf-example"
  description     = "tf-example-description"
  protocol        = "udp"
  to_destination  = "172.17.20.30"
}
```

## Argument Reference
- `dc_interface_id` **required** *string* &rarr;  Direct Connect Interface ID. Changing this creates a new resource

- `protocol` **required** *string* &rarr;  Protocol. Must be one of: "tcp", "udp", "any".

- `to_destination` **required** *string* &rarr;  IP Address of forwarding's destination.

- `description` optional *string* &rarr;  Description of the conntrack helper

- `destination` optional *string* &rarr;  Destination address selector.

- `name` optional *string* &rarr;  Name of the conntrack helper

- `port` optional *number* &rarr;  Port selector.

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.

- `source` optional *string* &rarr;  Source address selector.

- `to_port` optional *number* &rarr;  Destination port selector.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `id` *string* &rarr;  ID of the resource

- `updated_at` *string* &rarr;  Update timestamp



## Import

Direct connect IP port forwarding can be imported using the `id`, e.g.
```shell
terraform import vkcs_dc_ip_port_forwarding.mydcipportforwarding 659be09e-a10e-4762-b729-7a003af40f29
```
