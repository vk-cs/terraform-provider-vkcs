---
subcategory: "Load Balancers"
layout: "vkcs"
page_title: "vkcs: vkcs_lb_monitor"
description: |-
  Manages a monitor resource within VKCS.
---

# vkcs_lb_monitor

Manages a monitor resource within VKCS.

## Example Usage
```terraform
resource "vkcs_lb_monitor" "monitor_1" {
	pool_id     = "${vkcs_lb_pool.pool_1.id}"
	type        = "PING"
	delay       = 20
	timeout     = 10
	max_retries = 5
}
```
## Argument Reference
- `delay` **required** *number* &rarr;  The time, in seconds, between sending probes to members.

- `max_retries` **required** *number* &rarr;  Number of permissible ping failures before changing the member's status to INACTIVE. Must be a number between 1 and 10.

- `pool_id` **required** *string* &rarr;  The id of the pool that this monitor will be assigned to.

- `timeout` **required** *number* &rarr;  Maximum number of seconds for a monitor to wait for a ping reply before it times out. The value must be less than the delay value.

- `type` **required** *string* &rarr;  The type of probe, which is PING, TCP, HTTP, HTTPS, TLS-HELLO or UDP-CONNECT, that is sent by the load balancer to verify the member state. Changing this creates a new monitor.

- `admin_state_up` optional *boolean* &rarr;  The administrative state of the monitor. A valid value is true (UP) or false (DOWN).

- `expected_codes` optional *string* &rarr;  Required for HTTP(S) types. Expected HTTP codes for a passing HTTP(S) monitor. You can either specify a single status like "200", or a range like "200-202".

- `http_method` optional *string* &rarr;  Required for HTTP(S) types. The HTTP method used for requests by the monitor. If this attribute is not specified, it defaults to "GET".

- `max_retries_down` optional *number* &rarr;  Number of permissible ping failures befor changing the member's status to ERROR. Must be a number between 1 and 10. Changing this updates the max_retries_down of the existing monitor.

- `name` optional *string* &rarr;  The Name of the Monitor.

- `region` optional *string* &rarr;  The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new monitor.

- `url_path` optional *string* &rarr;  Required for HTTP(S) types. URI path that will be accessed if monitor type is HTTP or HTTPS.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Load Balancer Pool Monitor can be imported using the Monitor ID, e.g.:

```shell
terraform import vkcs_lb_monitor.monitor_1 47c26fc3-2403-427a-8c79-1589bd0533c2
```
