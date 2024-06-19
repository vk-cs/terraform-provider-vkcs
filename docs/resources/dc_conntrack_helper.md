---
subcategory: "Direct Connect"
layout: "vkcs"
page_title: "vkcs: vkcs_dc_conntrack_helper"
description: |-
  Manages a direct connect conntrack helper resource within VKCS.
---

# vkcs_dc_conntrack_helper

Manages a direct connect conntrack helper resource.

~> **Note:** This resource requires Sprut SDN to be enabled in your project.

**New since v0.8.0**.

## Example Usage
```terraform
resource "vkcs_dc_conntrack_helper" "dc-conntrack-helper" {
  dc_router_id = vkcs_dc_router.dc_router.id
  name         = "tf-example"
  description  = "tf-example-description"
  helper       = "ftp"
  protocol     = "tcp"
  port         = 21
}
```

## Argument Reference
- `dc_router_id` **required** *string* &rarr;  Direct Connect Router ID. Changing this creates a new resource

- `helper` **required** *string* &rarr;  Helper type. Must be one of: "ftp".

- `port` **required** *number* &rarr;  Network port for conntrack target rule.

- `protocol` **required** *string* &rarr;  Protocol. Must be one of: "tcp".

- `description` optional *string* &rarr;  Description of the conntrack helper

- `name` optional *string* &rarr;  Name of the conntrack helper

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `id` *string* &rarr;  ID of the resource

- `updated_at` *string* &rarr;  Update timestamp



## Import

Direct connect conntrack helper can be imported using the `id`, e.g.
```shell
terraform import vkcs_dc_conntrack_helper.mydcconntrackhelper 52a0a638-0a75-4a15-b3f3-d5c9f953e93f
```
