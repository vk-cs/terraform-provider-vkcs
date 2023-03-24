---
layout: "vkcs"
page_title: "vkcs: vkcs_lb_loadbalancer"
description: |-
  Get information on a VKCS Loadbalancer
---

# vkcs_lb_loadbalancer

Use this data source to get the details of a loadbalancer

## Example Usage

```terraform
data "vkcs_lb_loadbalancer" "loadbalancer" {
  id = "35082f6e-14c4-478c-ba4c-77bcdb222743"
}
```

## Argument Reference
- `id` **String** (***Required***) The UUID of the Loadbalancer

- `region` **String** (*Optional*) The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
- `id` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `admin_state_up` **Boolean** The administrative state of the Loadbalancer.

- `availability_zone` **String** The availability zone of the Loadbalancer.

- `description` **String** Human-readable description of the Loadbalancer.

- `name` **String** The name of the Loadbalancer.

- `security_group_ids` <strong>Set of </strong>**String** A list of security group IDs applied to the Loadbalancer.

- `tags` <strong>Set of </strong>**String** A list of simple strings assigned to the loadbalancer.

- `vip_address` **String** The ip address of the Loadbalancer.

- `vip_network_id` **String** The network on which to allocate the Loadbalancer's address. A tenant can only create Loadbalancers on networks authorized by policy (e.g. networks that belong to them or networks that are shared).  Changing this creates a new loadbalancer.

- `vip_port_id` **String** The port UUID of the Loadbalancer.

- `vip_subnet_id` **String** The subnet on which the Loadbalancer's address is allocated.


