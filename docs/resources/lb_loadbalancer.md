---
subcategory: "Load Balancers"
layout: "vkcs"
page_title: "vkcs: vkcs_lb_loadbalancer"
description: |-
  Manages a loadbalancer resource within VKCS.
---

# vkcs_lb_loadbalancer

Manages a loadbalancer resource within VKCS.

## Example Usage
```terraform
resource "vkcs_lb_loadbalancer" "lb_1" {
	vip_subnet_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"
}
```
## Argument Reference
- `admin_state_up` optional *boolean* &rarr;  The administrative state of the Loadbalancer. A valid value is true (UP) or false (DOWN).

- `availability_zone` optional *string* &rarr;  The availability zone of the Loadbalancer. Changing this creates a new loadbalancer.

- `description` optional *string* &rarr;  Human-readable description for the Loadbalancer.

- `name` optional *string* &rarr;  Human-readable name for the Loadbalancer. Does not have to be unique.

- `region` optional *string* &rarr;  The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new LB loadbalancer.

- `security_group_ids` optional deprecated *set of* *string* &rarr;  A list of security group IDs to apply to the loadbalancer. The security groups must be specified by ID and not name (as opposed to how they are configured with the Compute Instance). **Deprecated** This argument is deprecated, please do not use it.

- `tags` optional *set of* *string* &rarr;  A list of simple strings assigned to the loadbalancer.

- `vip_address` optional *string* &rarr;  The ip address of the load balancer. Changing this creates a new loadbalancer.

- `vip_network_id` optional *string* &rarr;  The network on which to allocate the Loadbalancer's address. A tenant can only create Loadbalancers on networks authorized by policy (e.g. networks that belong to them or networks that are shared).  Changing this creates a new loadbalancer.

- `vip_port_id` optional *string* &rarr;  The port UUID that the loadbalancer will use. Changing this creates a new loadbalancer.

- `vip_subnet_id` optional *string* &rarr;  The subnet on which to allocate the Loadbalancer's address. A tenant can only create Loadbalancers on networks authorized by policy (e.g. networks that belong to them or networks that are shared).  Changing this creates a new loadbalancer.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Load Balancer can be imported using the Load Balancer ID, e.g.:

```shell
terraform import vkcs_lb_loadbalancer.loadbalancer_1 19bcfdc7-c521-4a7e-9459-6750bd16df76
```
