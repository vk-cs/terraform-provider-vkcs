---
subcategory: "Load Balancers"
layout: "vkcs"
page_title: "vkcs: vkcs_lb_l7policy"
description: |-
  Manages a L7 Policy resource within VKCS.
---

# vkcs_lb_l7policy

Manages a Load Balancer L7 Policy resource within VKCS.

## Example Usage
```terraform
resource "vkcs_lb_l7policy" "app_redirect" {
  name             = "http-tf-example"
  description      = "Policy for tf lb testing"
  action           = "REDIRECT_TO_POOL"
  position         = 1
  listener_id      = vkcs_lb_listener.app_http.id
  redirect_pool_id = vkcs_lb_pool.http.id
}
```
## Argument Reference
- `action` **required** *string* &rarr;  The L7 Policy action - can either be REDIRECT\_TO\_POOL, REDIRECT\_TO\_URL or REJECT.

- `listener_id` **required** *string* &rarr;  The Listener on which the L7 Policy will be associated with. Changing this creates a new L7 Policy.

- `admin_state_up` optional *boolean* &rarr;  The administrative state of the L7 Policy. A valid value is true (UP) or false (DOWN).

- `description` optional *string* &rarr;  Human-readable description for the L7 Policy.

- `name` optional *string* &rarr;  Human-readable name for the L7 Policy. Does not have to be unique.

- `position` optional *number* &rarr;  The position of this policy on the listener. Positions start at 1.

- `redirect_pool_id` optional *string* &rarr;  Requests matching this policy will be redirected to the pool with this ID. Only valid if action is REDIRECT\_TO\_POOL.

- `redirect_url` optional *string* &rarr;  Requests matching this policy will be redirected to this URL. Only valid if action is REDIRECT\_TO\_URL.

- `region` optional *string* &rarr;  The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new L7 Policy.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Load Balancer L7 Policy can be imported using the L7 Policy ID, e.g.:

```shell
terraform import vkcs_lb_l7policy.l7policy_1 8a7a79c2-cf17-4e65-b2ae-ddc8bfcf6c74
```
