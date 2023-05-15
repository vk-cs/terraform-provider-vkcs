---
subcategory: "Load Balancers"
layout: "vkcs"
page_title: "vkcs: vkcs_lb_members"
description: |-
  Manages a members resource within VKCS.
---

# vkcs_lb_members



## Example Usage
```terraform
resource "vkcs_lb_members" "members_1" {
	pool_id = "935685fb-a896-40f9-9ff4-ae531a3a00fe"

	member {
		address       = "192.168.199.23"
		protocol_port = 8080
	}

	member {
		address       = "192.168.199.24"
		protocol_port = 8080
	}
}
```
## Argument Reference
- `pool_id` **required** *string* &rarr;  The id of the pool that members will be assigned to. Changing this creates a new members resource.

- `member` optional &rarr;  A set of dictionaries containing member parameters. The structure is described below.
  - `address` **required** *string* &rarr;  The IP address of the members to receive traffic from the load balancer.

  - `protocol_port` **required** *number* &rarr;  The port on which to listen for client traffic.

  - `admin_state_up` optional *boolean* &rarr;  The administrative state of the member. A valid value is true (UP) or false (DOWN). Defaults to true.

  - `backup` optional *boolean* &rarr;  A bool that indicates whether the member is backup.

  - `name` optional *string* &rarr;  Human-readable name for the member.

  - `subnet_id` optional *string* &rarr;  The subnet in which to access the member.

  - `weight` optional *number* &rarr;  A positive integer value that indicates the relative portion of traffic that this members should receive from the pool. For example, a member with a weight of 10 receives five times as much traffic as a member with a weight of 2. Defaults to 1.

- `region` optional *string* &rarr;  The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new members resource.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `member` 
  - `id` *string* &rarr;  The unique ID for the member.



## Import

Load Balancer Pool Members can be imported using the Pool ID, e.g.:

```shell
terraform import vkcs_lb_members.members_1 c22974d2-4c95-4bcb-9819-0afc5ed303d5
```
