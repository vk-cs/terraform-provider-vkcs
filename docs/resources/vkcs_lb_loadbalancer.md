---
layout: "vkcs"
page_title: "vkcs: lb_loadbalancer"
description: |-
	Manages a loadbalancer resource within VKCS.
---

# vkcs\_lb\_loadbalancer

Manages a loadbalancer resource within VKCS.

## Example Usage

```hcl
resource "vkcs_lb_loadbalancer" "lb_1" {
	vip_subnet_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"
}
```

## Argument Reference

The following arguments are supported:

* `admin_state_up` - (Optional) The administrative state of the Loadbalancer.
	A valid value is true (UP) or false (DOWN).

* `availability_zone` - (Optional) The availability zone of the Loadbalancer.
  Changing this creates a new loadbalancer.

* `description` - (Optional) Human-readable description for the Loadbalancer.

* `region` - (Optional) The region in which to obtain the Loadbalancer client.
	If omitted, the	`region` argument of the provider is used. Changing this creates a new
	LB loadbalancer.

* `name` - (Optional) Human-readable name for the Loadbalancer. Does not have
	to be unique.

* `security_group_ids` - (Optional) A list of security group IDs to apply to the
	loadbalancer. The security groups must be specified by ID and not name (as
	opposed to how they are configured with the Compute Instance).

* `tags` - (Optional) A list of simple strings assigned to the loadbalancer.

* `vip_address` - (Optional) The ip address of the load balancer.
	Changing this creates a new loadbalancer.

* `vip_network_id` - (Optional) The network on which to allocate the
	Loadbalancer's address. A tenant can only create Loadbalancers on networks
	authorized by policy (e.g. networks that belong to them or networks that
	are shared).  Changing this creates a new loadbalancer.

* `vip_port_id` - (Optional) The port UUID that the loadbalancer will use.
  Changing this creates a new loadbalancer.

* `vip_subnet_id` - (Optional) The subnet on which to allocate the
	Loadbalancer's address. A tenant can only create Loadbalancers on networks
	authorized by policy (e.g. networks that belong to them or networks that
	are shared).  Changing this creates a new loadbalancer.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `vip_subnet_id` - See Argument Reference above.
* `vip_network_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `availability_zone` - See Argument Reference above.
* `security_group_ids` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `vip_port_id` - The Port ID of the Load Balancer IP.
* `vip_address` - See Argument Reference above.

## Import

Load Balancer can be imported using the Load Balancer ID, e.g.:

```
$ terraform import vkcs_lb_loadbalancer.loadbalancer_1 19bcfdc7-c521-4a7e-9459-6750bd16df76
```
