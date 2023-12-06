---
subcategory: "Load Balancers"
layout: "vkcs"
page_title: "vkcs: vkcs_lb_loadbalancer"
description: |-
  Get information on a VKCS Loadbalancer
---

# vkcs_lb_loadbalancer

Use this data source to get the details of a loadbalancer

## Example Usage

```terraform
data "vkcs_lb_loadbalancer" "app" {
  id = vkcs_lb_loadbalancer.app.id
  # This is unnecessary in real life.
  # This is required here to let the example work with loadbalancer resource example. 
  depends_on = [vkcs_lb_loadbalancer.app]
}

data "vkcs_networking_port" "app-port" {
  port_id = data.vkcs_lb_loadbalancer.app.vip_port_id
}

output "used_vips" {
  value = data.vkcs_networking_port.app-port.all_fixed_ips
  description = "IP addresses of the app"
}
```

## Argument Reference
- `id` **required** *string* &rarr;  The UUID of the Loadbalancer

- `region` optional *string* &rarr;  The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `admin_state_up` *boolean* &rarr;  The administrative state of the Loadbalancer.

- `availability_zone` *string* &rarr;  The availability zone of the Loadbalancer.

- `description` *string* &rarr;  Human-readable description of the Loadbalancer.

- `name` *string* &rarr;  The name of the Loadbalancer.

- `security_group_ids` *set of* *string* &rarr;  A list of security group IDs applied to the Loadbalancer.

- `tags` *set of* *string* &rarr;  A list of simple strings assigned to the loadbalancer.

- `vip_address` *string* &rarr;  The ip address of the Loadbalancer.

- `vip_network_id` *string* &rarr;  The network on which to allocate the Loadbalancer's address. A tenant can only create Loadbalancers on networks authorized by policy (e.g. networks that belong to them or networks that are shared).  Changing this creates a new loadbalancer.

- `vip_port_id` *string* &rarr;  The port UUID of the Loadbalancer.

- `vip_subnet_id` *string* &rarr;  The subnet on which the Loadbalancer's address is allocated.


