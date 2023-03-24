---
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
- `admin_state_up` **Boolean** (*Optional*) The administrative state of the Loadbalancer. A valid value is true (UP) or false (DOWN).

- `availability_zone` **String** (*Optional*) The availability zone of the Loadbalancer. Changing this creates a new loadbalancer.

- `description` **String** (*Optional*) Human-readable description for the Loadbalancer.

- `name` **String** (*Optional*) Human-readable name for the Loadbalancer. Does not have to be unique.

- `region` **String** (*Optional*) The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new LB loadbalancer.

- `security_group_ids` <strong>Set of </strong>**String** (*Optional* Deprecated) A list of security group IDs to apply to the loadbalancer. The security groups must be specified by ID and not name (as opposed to how they are configured with the Compute Instance). ***Deprecated*** This argument is deprecated, please do not use it.

- `tags` <strong>Set of </strong>**String** (*Optional*) A list of simple strings assigned to the loadbalancer.

- `vip_address` **String** (*Optional*) The ip address of the load balancer. Changing this creates a new loadbalancer.

- `vip_network_id` **String** (*Optional*) The network on which to allocate the Loadbalancer's address. A tenant can only create Loadbalancers on networks authorized by policy (e.g. networks that belong to them or networks that are shared).  Changing this creates a new loadbalancer.

- `vip_port_id` **String** (*Optional*) The port UUID that the loadbalancer will use. Changing this creates a new loadbalancer.

- `vip_subnet_id` **String** (*Optional*) The subnet on which to allocate the Loadbalancer's address. A tenant can only create Loadbalancers on networks authorized by policy (e.g. networks that belong to them or networks that are shared).  Changing this creates a new loadbalancer.


## Attributes Reference
- `admin_state_up` **Boolean** See Argument Reference above.

- `availability_zone` **String** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `security_group_ids` <strong>Set of </strong>**String** See Argument Reference above.

- `tags` <strong>Set of </strong>**String** See Argument Reference above.

- `vip_address` **String** See Argument Reference above.

- `vip_network_id` **String** See Argument Reference above.

- `vip_port_id` **String** See Argument Reference above.

- `vip_subnet_id` **String** See Argument Reference above.

- `id` **String** ID of the resource.



## Import

Load Balancer can be imported using the Load Balancer ID, e.g.:

```shell
terraform import vkcs_lb_loadbalancer.loadbalancer_1 19bcfdc7-c521-4a7e-9459-6750bd16df76
```
