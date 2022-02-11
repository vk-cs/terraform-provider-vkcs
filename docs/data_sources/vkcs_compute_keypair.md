---
layout: "vkcs"
page_title: "vkcs: compute_keypair"
description: |-
  Get information on an OpenStack Keypair.
---

# vkcs\_compute\_keypair

Use this data source to get the ID and public key of an OpenStack keypair.

## Example Usage

```hcl
data "vkcs_compute_keypair" "kp" {
  name = "sand"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    If omitted, the `region` argument of the provider is used.

* `name` - (Required) The unique name of the keypair.


## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `fingerprint` - The fingerprint of the OpenSSH key.
* `public_key` - The OpenSSH-formatted public key of the keypair.
