---
subcategory: "VPN"
layout: "vkcs"
page_title: "vkcs: vkcs_vpnaas_ike_policy"
description: |-
  Manages a IKE policy resource within VKCS.
---

# vkcs_vpnaas_ike_policy

Manages a IKE policy resource within VKCS.

## Example Usage
```terraform
resource "vkcs_vpnaas_ike_policy" "policy_2" {
	name = "my_policy"
}
```
## Argument Reference
- `auth_algorithm` optional *string* &rarr;  The authentication hash algorithm. Valid values are sha1, sha256, sha384, sha512. Default is sha1. Changing this updates the algorithm of the existing policy.

- `description` optional *string* &rarr;  The human-readable description for the policy. Changing this updates the description of the existing policy.

- `encryption_algorithm` optional *string* &rarr;  The encryption algorithm. Valid values are 3des, aes-128, aes-192 and so on. The default value is aes-128. Changing this updates the existing policy.

- `ike_version` optional *string* &rarr;  The IKE mode. A valid value is v1 or v2. Default is v1. Changing this updates the existing policy.

- `lifetime` optional &rarr;  The lifetime of the security association. Consists of Unit and Value.
  - `units` optional *string* &rarr;  The units for the lifetime of the security association. Can be either seconds or kilobytes. Default is seconds.

  - `value` optional *number* &rarr;  The value for the lifetime of the security association. Must be a positive integer. Default is 3600.

- `name` optional *string* &rarr;  The name of the policy. Changing this updates the name of the existing policy.

- `pfs` optional *string* &rarr;  The perfect forward secrecy mode. Valid values are Group2, Group5 and Group14. Default is Group5. Changing this updates the existing policy.

- `phase1_negotiation_mode` optional *string* &rarr;  The IKE mode. A valid value is main, which is the default. Changing this updates the existing policy.

- `region` optional *string* &rarr;  The region in which to obtain the Networking client. A Networking client is needed to create a VPN service. If omitted, the `region` argument of the provider is used. Changing this creates a new service.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Services can be imported using the `id`, e.g.

```shell
terraform import vkcs_vpnaas_ike_policy.policy_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
