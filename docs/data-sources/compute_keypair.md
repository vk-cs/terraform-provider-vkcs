---
subcategory: "Virtual Machines"
layout: "vkcs"
page_title: "vkcs: vkcs_compute_keypair"
description: |-
  Get information on an VKCS Keypair.
---

# vkcs_compute_keypair

Use this data source to get the ID and public key of an VKCS keypair.

## Example Usage

```terraform
data "vkcs_compute_keypair" "generated_key" {
  name = "generated-key-tf-example"
  # This is unnecessary in real life.
  # This is required here to let the example work with keypair resource example. 
  depends_on = [vkcs_compute_keypair.generated_key]
}
```

## Argument Reference
- `name` **required** *string* &rarr;  The unique name of the keypair.

- `region` optional *string* &rarr;  The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `fingerprint` *string* &rarr;  The fingerprint of the OpenSSH key.

- `id` *string* &rarr;  ID of the resource.

- `public_key` *string* &rarr;  The OpenSSH-formatted public key of the keypair.


