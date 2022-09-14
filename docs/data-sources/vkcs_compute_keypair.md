---
layout: "vkcs"
page_title: "vkcs: vkcs_compute_keypair"
description: |-
  Get information on an VKCS Keypair.
---

# vkcs_compute_keypair

Use this data source to get the ID and public key of an VKCS keypair.

## Example Usage

```terraform
data "vkcs_compute_keypair" "kp" {
  name = "sand"
}
```

## Argument Reference
- `name` **String** (***Required***) The unique name of the keypair.

- `region` **String** (*Optional*) The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
- `name` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `fingerprint` **String** The fingerprint of the OpenSSH key.

- `id` **String** ID of the resource.

- `public_key` **String** The OpenSSH-formatted public key of the keypair.


