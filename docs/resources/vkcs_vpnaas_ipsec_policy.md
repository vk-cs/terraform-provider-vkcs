---
layout: "vkcs"
page_title: "vkcs: vpnaas_ipsec_policy"
description: |-
	Manages a IPSec policy resource within VKCS.
---

# vkcs\_vpnaas\_ipsec\_policy

Manages a IPSec policy resource within VKCS.

## Example Usage

```hcl
resource "vkcs_vpnaas_ipsec_policy" "policy_1" {
	name = "my_policy"
}
```

## Argument Reference

The following arguments are supported:

* `auth_algorithm` - (Optional) The authentication hash algorithm. Valid values are sha1, sha256, sha384, sha512.
	Default is sha1. Changing this updates the algorithm of the existing policy.

* `description` - (Optional) The human-readable description for the policy.
	Changing this updates the description of the existing policy.

* `encapsulation_mode` - (Optional) The encapsulation mode. Valid values are tunnel and transport. Default is tunnel.
	Changing this updates the existing policy.

* `encryption_algorithm` - (Optional) The encryption algorithm. Valid values are 3des, aes-128, aes-192 and so on.
	The default value is aes-128. Changing this updates the existing policy.

* `lifetime` - (Optional) The lifetime of the security association. Consists of Unit and Value.
	- `unit` - (Optional) The units for the lifetime of the security association. Can be either seconds or kilobytes.
	Default is seconds.
	- `value` - (Optional) The value for the lifetime of the security association. Must be a positive integer.
	Default is 3600.

* `name` - (Optional) The name of the policy. Changing this updates the name of
	the existing policy.

* `pfs` - (Optional) The perfect forward secrecy mode. Valid values are Group2, Group5 and Group14. Default is Group5.
	Changing this updates the existing policy.

* `region` - (Optional) The region in which to obtain the Networking client.
	A Networking client is needed to create an IPSec policy. If omitted, the
	`region` argument of the provider is used. Changing this creates a new
	policy.

* `transform_protocol` - (Optional) The transform protocol. Valid values are ESP, AH and AH-ESP.
	Changing this updates the existing policy. Default is ESP.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `auth_algorithm` - See Argument Reference above.
* `encapsulation_mode` - See Argument Reference above.
* `encryption_algorithm` - See Argument Reference above.
* `pfs` - See Argument Reference above.
* `transform_protocol` - See Argument Reference above.
* `lifetime` - See Argument Reference above.
	- `unit` - See Argument Reference above.
	- `value` - See Argument Reference above.

## Import

Policies can be imported using the `id`, e.g.

```
$ terraform import vkcs_vpnaas_ipsec_policy.policy_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
