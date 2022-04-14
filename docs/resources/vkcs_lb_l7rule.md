---
layout: "vkcs"
page_title: "vkcs: lb_l7rule"
description: |-
	Manages a L7 rule resource within VKCS.
---

# vkcs\_lb\_l7rule

Manages a L7 Rule resource within VKCS.

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
	name         = "test"
	action       = "REDIRECT_TO_URL"
	description  = "test description"
	position     = 1
	listener_id  = "${vkcs_lb_listener.listener_1.id}"
	redirect_url = "http://www.example.com"
}

resource "vkcs_lb_l7rule" "l7rule_1" {
	l7policy_id  = "${vkcs_lb_l7policy.l7policy_1.id}"
	compare_type = "EQUAL_TO"
	type         = "PATH"
	value        = "/api"
}
```

## Argument Reference

The following arguments are supported:

* `l7policy_id` - (Required) The ID of the L7 Policy to query. Changing this creates a new
	L7 Rule.

* `compare_type` - (Required) The comparison type for the L7 rule - can either be
	CONTAINS, STARTS\_WITH, ENDS_WITH, EQUAL_TO or REGEX

* `type` - (Required) The L7 Rule type - can either be COOKIE, FILE\_TYPE, HEADER,
	HOST\_NAME or PATH.

* `value` - (Required) The value to use for the comparison. For example, the file type to
	compare.

* `admin_state_up` - (Optional) The administrative state of the L7 Rule.
	A valid value is true (UP) or false (DOWN).

* `description` - (Optional) Human-readable description for the L7 Rule.

* `invert` - (Optional) When true the logic of the rule is inverted. For example, with invert
	true, equal to would become not equal to. Default is false.

* `key` - (Optional) The key to use for the comparison. For example, the name of the cookie to
	evaluate. Valid when `type` is set to COOKIE or HEADER.
* `region` - (Optional) The region in which to obtain the Loadbalancer client.
	If omitted, the `region` argument of the provider is used. Changing this creates a new
	L7 Rule.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the L7 Rule.
* `region` - See Argument Reference above.
* `type` - See Argument Reference above.
* `compare_type` - See Argument Reference above.
* `l7policy_id` - See Argument Reference above.
* `value` - See Argument Reference above.
* `key` - See Argument Reference above.
* `invert` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `listener_id` - The ID of the Listener owning this resource.

## Import

Load Balancer L7 Rule can be imported using the L7 Policy ID and L7 Rule ID
separated by a slash, e.g.:

```
$ terraform import vkcs_lb_l7rule.l7rule_1 e0bd694a-abbe-450e-b329-0931fd1cc5eb/4086b0c9-b18c-4d1c-b6b8-4c56c3ad2a9e
```
