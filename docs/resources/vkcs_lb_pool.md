---
layout: "vkcs"
page_title: "vkcs: vkcs_lb_pool"
description: |-
  Manages a pool resource within VKCS.
---

# vkcs_lb_pool

Manages a pool resource within VKCS.

## Example Usage
```terraform
resource "vkcs_lb_pool" "pool_1" {
	protocol    = "HTTP"
	lb_method   = "ROUND_ROBIN"
	listener_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"

	persistence {
		type        = "APP_COOKIE"
		cookie_name = "testCookie"
	}
}
```
## Argument Reference
- `lb_method` **String** (***Required***) The load balancing algorithm to distribute traffic to the pool's members. Must be one of ROUND_ROBIN, LEAST_CONNECTIONS, SOURCE_IP, or SOURCE_IP_PORT.

- `protocol` **String** (***Required***) The protocol - can either be TCP, HTTP, HTTPS, PROXY, or UDP. Changing this creates a new pool.

- `admin_state_up` **Boolean** (*Optional*) The administrative state of the pool. A valid value is true (UP) or false (DOWN).

- `description` **String** (*Optional*) Human-readable description for the pool.

- `listener_id` **String** (*Optional*) The Listener on which the members of the pool will be associated with. Changing this creates a new pool. Note:  One of LoadbalancerID or ListenerID must be provided.

- `loadbalancer_id` **String** (*Optional*) The load balancer on which to provision this pool. Changing this creates a new pool. Note: One of LoadbalancerID or ListenerID must be provided.

- `name` **String** (*Optional*) Human-readable name for the pool.

- `persistence` (*Optional*) Omit this field to prevent session persistence. Indicates whether connections in the same session will be processed by the same Pool member or not. Changing this creates a new pool.
  - `type` **String** (***Required***) The type of persistence mode. The current specification supports SOURCE_IP, HTTP_COOKIE, and APP_COOKIE.

  - `cookie_name` **String** (*Optional*) The name of the cookie if persistence mode is set appropriately. Required if `type = APP_COOKIE`.

- `region` **String** (*Optional*) The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new pool.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` **String** ID of the resource.



## Import

Load Balancer Pool can be imported using the Pool ID, e.g.:

```shell
terraform import vkcs_lb_pool.pool_1 60ad9ee4-249a-4d60-a45b-aa60e046c513
```
