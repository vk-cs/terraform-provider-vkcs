---
layout: "vkcs"
page_title: "VKCS: lb_l7policy"
description: |-
	Manages a L7 Policy resource within OpenStack.
---

# vkcs\_lb\_l7policy

Manages a Load Balancer L7 Policy resource within OpenStack.

## Example Usage

```hcl
resource "vkcs_networking_network" "network_1" {
	name           = "network_1"
	admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
	name       = "subnet_1"
	cidr       = "192.168.199.0/24"
	ip_version = 4
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

The following arguments are supported:

* `action` - (Required) The L7 Policy action - can either be REDIRECT\_TO\_POOL,
	REDIRECT\_TO\_URL or REJECT.

* `listener_id` - (Required) The Listener on which the L7 Policy will be associated with.
	Changing this creates a new L7 Policy.

* `admin_state_up` - (Optional) The administrative state of the L7 Policy.
	A valid value is true (UP) or false (DOWN).

* `description` - (Optional) Human-readable description for the L7 Policy.

* `name` - (Optional) Human-readable name for the L7 Policy. Does not have
	to be unique.

* `position` - (Optional) The position of this policy on the listener. Positions start at 1.

* `redirect_pool_id` - (Optional) Requests matching this policy will be redirected to the
	pool with this ID. Only valid if action is REDIRECT\_TO\_POOL.

* `redirect_url` - (Optional) Requests matching this policy will be redirected to this URL.
	Only valid if action is REDIRECT\_TO\_URL.

* `region` - (Optional) The region in which to obtain the V2 Networking client.
	A Networking client is needed to create an . If omitted, the
	`region` argument of the provider is used. Changing this creates a new
	L7 Policy.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the L7 Policy.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `action` - See Argument Reference above.
* `listener_id` - See Argument Reference above.
* `position` - See Argument Reference above.
* `redirect_pool_id` - See Argument Reference above.
* `redirect_url` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.

## Import

Load Balancer L7 Policy can be imported using the L7 Policy ID, e.g.:

```
$ terraform import vkcs_lb_l7policy.l7policy_1 8a7a79c2-cf17-4e65-b2ae-ddc8bfcf6c74
```
