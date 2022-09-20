---
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
- `auth_algorithm` **String** (*Optional*) The authentication hash algorithm. Valid values are sha1, sha256, sha384, sha512. Default is sha1. Changing this updates the algorithm of the existing policy.

- `description` **String** (*Optional*) The human-readable description for the policy. Changing this updates the description of the existing policy.

- `encryption_algorithm` **String** (*Optional*) The encryption algorithm. Valid values are 3des, aes-128, aes-192 and so on. The default value is aes-128. Changing this updates the existing policy.

- `ike_version` **String** (*Optional*) The IKE mode. A valid value is v1 or v2. Default is v1. Changing this updates the existing policy.

- `lifetime` (*Optional*) The lifetime of the security association. Consists of Unit and Value.
  - `units` **String** (*Optional*) The units for the lifetime of the security association. Can be either seconds or kilobytes. Default is seconds.

  - `value` **Number** (*Optional*) The value for the lifetime of the security association. Must be a positive integer. Default is 3600.

- `name` **String** (*Optional*) The name of the policy. Changing this updates the name of the existing policy.

- `pfs` **String** (*Optional*) The perfect forward secrecy mode. Valid values are Group2, Group5 and Group14. Default is Group5. Changing this updates the existing policy.

- `phase1_negotiation_mode` **String** (*Optional*) The IKE mode. A valid value is main, which is the default. Changing this updates the existing policy.

- `region` **String** (*Optional*) The region in which to obtain the Networking client. A Networking client is needed to create a VPN service. If omitted, the `region` argument of the provider is used. Changing this creates a new service.


## Attributes Reference
- `auth_algorithm` **String** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `encryption_algorithm` **String** See Argument Reference above.

- `ike_version` **String** See Argument Reference above.

- `lifetime`  See Argument Reference above.
  - `units` **String** See Argument Reference above.

  - `value` **Number** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `pfs` **String** See Argument Reference above.

- `phase1_negotiation_mode` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `id` **String** ID of the resource.



## Import

Services can be imported using the `id`, e.g.

```shell
terraform import vkcs_vpnaas_ike_policy.policy_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
