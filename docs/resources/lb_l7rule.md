---
subcategory: "Load Balancers"
layout: "vkcs"
page_title: "vkcs: vkcs_lb_l7rule"
description: |-
  Manages a L7 rule resource within VKCS.
---

# vkcs_lb_l7rule

Manages a L7 Rule resource within VKCS.

## Example Usage
```terraform
resource "vkcs_lb_l7rule" "app_api_redirect" {
  l7policy_id  = vkcs_lb_l7policy.app_redirect.id
  compare_type = "EQUAL_TO"
  type         = "PATH"
  value        = "/api"
}
```
## Argument Reference
- `compare_type` **required** *string* &rarr;  The comparison type for the L7 rule - can either be CONTAINS, STARTS\_WITH, ENDS_WITH, EQUAL_TO or REGEX

- `l7policy_id` **required** *string* &rarr;  The ID of the L7 Policy to query. Changing this creates a new L7 Rule.

- `type` **required** *string* &rarr;  The L7 Rule type - can either be COOKIE, FILE\_TYPE, HEADER, HOST\_NAME or PATH.

- `value` **required** *string* &rarr;  The value to use for the comparison. For example, the file type to compare.

- `admin_state_up` optional *boolean* &rarr;  The administrative state of the L7 Rule. A valid value is true (UP) or false (DOWN).

- `invert` optional *boolean* &rarr;  When true the logic of the rule is inverted. For example, with invert true, equal to would become not equal to. Default is false.

- `key` optional *string* &rarr;  The key to use for the comparison. For example, the name of the cookie to evaluate. Valid when `type` is set to COOKIE or HEADER.

- `region` optional *string* &rarr;  The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new L7 Rule.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `listener_id` *string* &rarr;  The ID of the Listener owning this resource.



## Import

Load Balancer L7 Rule can be imported using the L7 Policy ID and L7 Rule ID separated by a slash, e.g.:

```shell
terraform import vkcs_lb_l7rule.l7rule_1 e0bd694a-abbe-450e-b329-0931fd1cc5eb/4086b0c9-b18c-4d1c-b6b8-4c56c3ad2a9e
```
