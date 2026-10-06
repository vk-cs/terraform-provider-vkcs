---
subcategory: "Virtual Machines"
layout: "vkcs"
page_title: "vkcs: vkcs_compute_servergroup"
description: |-
  Get information on an VKCS server group.
---

# vkcs_compute_servergroup

Use this data source to get the ID of an available VKCS server group.

## Example Usage

```terraform
data "vkcs_compute_servergroup" "example" {
  name = "servergroup-tf-example"
}
```

## Argument Reference
- `name` **required** *string* &rarr;  The name of the server group.

- `region` optional *string* &rarr;  The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `members` *string* &rarr;  The instances that are part of this server group.

- `metadata` *map of* *string* &rarr;  The metadata of the server group.

- `policies` *string* &rarr;  The set of policies for the server group.

