---
subcategory: "VPN"
layout: "vkcs"
page_title: "vkcs: vkcs_vpnaas_ipsec_policy"
description: |-
  Manages a IPSec policy resource within VKCS.
---

# vkcs_vpnaas_ipsec_policy

Manages a IPSec policy resource within VKCS.

## Example Usage
```terraform
resource "vkcs_vpnaas_ipsec_policy" "policy_1" {
	name = "my_policy"
}
```
## Argument Reference
- `auth_algorithm` optional *string* &rarr;  The authentication hash algorithm. Valid values are sha1, sha256, sha384, sha512. Default is sha1. Changing this updates the algorithm of the existing policy.

- `description` optional *string* &rarr;  The human-readable description for the policy. Changing this updates the description of the existing policy.

- `encapsulation_mode` optional *string* &rarr;  The encapsulation mode. Valid values are tunnel and transport. Default is tunnel. Changing this updates the existing policy.

- `encryption_algorithm` optional *string* &rarr;  The encryption algorithm. Valid values are 3des, aes-128, aes-192 and so on. The default value is aes-128. Changing this updates the existing policy.

- `lifetime` optional &rarr;  The lifetime of the security association. Consists of Unit and Value.
  - `units` optional *string* &rarr;  The units for the lifetime of the security association. Can be either seconds or kilobytes. Default is seconds.

  - `value` optional *number* &rarr;  The value for the lifetime of the security association. Must be a positive integer. Default is 3600.

- `name` optional *string* &rarr;  The name of the policy. Changing this updates the name of the existing policy.

- `pfs` optional *string* &rarr;  The perfect forward secrecy mode. Valid values are Group2, Group5 and Group14. Default is Group5. Changing this updates the existing policy.

- `region` optional *string* &rarr;  The region in which to obtain the Networking client. A Networking client is needed to create an IPSec policy. If omitted, the `region` argument of the provider is used. Changing this creates a new policy.

- `transform_protocol` optional *string* &rarr;  The transform protocol. Valid values are ESP, AH and AH-ESP. Changing this updates the existing policy. Default is ESP.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Services can be imported using the `id`, e.g.

```shell
terraform import vkcs_vpnaas_ipsec_policy.policy_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
