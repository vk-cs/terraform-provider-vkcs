---
layout: "vkcs"
page_title: "vkcs: vkcs_lb_l7policy"
description: |-
  Manages a L7 Policy resource within VKCS.
---

# vkcs_lb_l7policy

Manages a Load Balancer L7 Policy resource within VKCS.

## Example Usage
```terraform
resource "vkcs_networking_network" "network_1" {
	name           = "network_1"
	admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
	name       = "subnet_1"
	cidr       = "192.168.199.0/24"
	network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
	name          = "loadbalancer_1"
	vip_subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}

resource "vkcs_lb_listener" "listener_1" {
	name            = "listener_1"
	protocol        = "HTTP"
	protocol_port   = 8080
	loadbalancer_id = "${vkcs_lb_loadbalancer.loadbalancer_1.id}"
}

resource "vkcs_lb_pool" "pool_1" {
	name            = "pool_1"
	protocol        = "HTTP"
	lb_method       = "ROUND_ROBIN"
	loadbalancer_id = "${vkcs_lb_loadbalancer.loadbalancer_1.id}"
}

resource "vkcs_lb_l7policy" "l7policy_1" {
	name             = "test"
	action           = "REDIRECT_TO_POOL"
	description      = "test l7 policy"
	position         = 1
	listener_id      = "${vkcs_lb_listener.listener_1.id}"
	redirect_pool_id = "${vkcs_lb_pool.pool_1.id}"
}
```
## Argument Reference
- `action` **String** (***Required***) The L7 Policy action - can either be REDIRECT\_TO\_POOL, REDIRECT\_TO\_URL or REJECT.

- `listener_id` **String** (***Required***) The Listener on which the L7 Policy will be associated with. Changing this creates a new L7 Policy.

- `admin_state_up` **Boolean** (*Optional*) The administrative state of the L7 Policy. A valid value is true (UP) or false (DOWN).

- `description` **String** (*Optional*) Human-readable description for the L7 Policy.

- `name` **String** (*Optional*) Human-readable name for the L7 Policy. Does not have to be unique.

- `position` **Number** (*Optional*) The position of this policy on the listener. Positions start at 1.

- `redirect_pool_id` **String** (*Optional*) Requests matching this policy will be redirected to the pool with this ID. Only valid if action is REDIRECT\_TO\_POOL.

- `redirect_url` **String** (*Optional*) Requests matching this policy will be redirected to this URL. Only valid if action is REDIRECT\_TO\_URL.

- `region` **String** (*Optional*) The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new L7 Policy.


## Attributes Reference
- `action` **String** See Argument Reference above.

- `listener_id` **String** See Argument Reference above.

- `admin_state_up` **Boolean** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `position` **Number** See Argument Reference above.

- `redirect_pool_id` **String** See Argument Reference above.

- `redirect_url` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `id` **String** ID of the resource.



## Import

Load Balancer L7 Policy can be imported using the L7 Policy ID, e.g.:

```shell
terraform import vkcs_lb_l7policy.l7policy_1 8a7a79c2-cf17-4e65-b2ae-ddc8bfcf6c74
```
