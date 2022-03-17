---
layout: "vkcs"
page_title: "VKCS: lb_listener"
description: |-
	Manages a listener resource within OpenStack.
---

# vkcs\_lb\_listener

Manages a listener resource within OpenStack.

~> **Note:** This resource has attributes that depend on octavia minor versions.
Please ensure your Openstack cloud supports the required [minor version](../#octavia-api-versioning).

## Example Usage

```hcl
resource "vkcs_lb_listener" "listener_1" {
	loadbalancer_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"
	protocol        = "HTTP"
	protocol_port   = 8080

	insert_headers = {
		X-Forwarded-For = "true"
	}
}
```

## Argument Reference

The following arguments are supported:

* `loadbalancer_id` - (Required) The load balancer on which to provision this
	Listener. Changing this creates a new Listener.

* `protocol` - (Required) The protocol - can either be TCP, HTTP, HTTPS,
	TERMINATED_HTTPS, UDP (supported only in Octavia) or SCTP (supported only
	in **Octavia minor version >= 2.23**). Changing this creates a new Listener.

* `protocol_port` - (Required) The port on which to listen for client traffic.
	Changing this creates a new Listener.

* `allowed_cidrs` - (Optional) A list of CIDR blocks that are permitted to connect to this listener, denying
	all other source addresses. If not present, defaults to allow all.

* `admin_state_up` - (Optional) The administrative state of the Listener.
	A valid value is true (UP) or false (DOWN).

* `connection_limit` - (Optional) The maximum number of connections allowed
	for the Listener.

* `default_pool_id` - (Optional) The ID of the default pool with which the
	Listener is associated.

* `default_tls_container_ref` - (Optional) A reference to a Barbican Secrets
	container which stores TLS information. This is required if the protocol
	is `TERMINATED_HTTPS`. See
	[here](https://wiki.openstack.org/wiki/Network/LBaaS/docs/how-to-create-tls-loadbalancer)
	for more information.

* `description` - (Optional) Human-readable description for the Listener.

* `insert_headers` - (Optional) The list of key value pairs representing headers to insert
	into the request before it is sent to the backend members. Changing this updates the headers of the
	existing listener.

* `name` - (Optional) Human-readable name for the Listener. Does not have
	to be unique.
* `region` - (Optional) The region in which to obtain the V2 Networking client.
	A Networking client is needed to create an . If omitted, the
	`region` argument of the provider is used. Changing this creates a new
	Listener.

* `sni_container_refs` - (Optional) A list of references to Barbican Secrets
	containers which store SNI information. See
	[here](https://wiki.openstack.org/wiki/Network/LBaaS/docs/how-to-create-tls-loadbalancer)
	for more information.

* `timeout_client_data` - (Optional) The client inactivity timeout in milliseconds.

* `timeout_member_connect` - (Optional) The member connection timeout in milliseconds.

* `timeout_member_data` - (Optional) The member inactivity timeout in milliseconds.

* `timeout_tcp_inspect` - (Optional) The time in milliseconds, to wait for additional
	TCP packets for content inspection.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the Listener.
* `protocol` - See Argument Reference above.
* `protocol_port` - See Argument Reference above.
* `name` - See Argument Reference above.
* `default_port_id` - See Argument Reference above.
* `description` - See Argument Reference above.
* `connection_limit` - See Argument Reference above.
* `timeout_client_data` - See Argument Reference above.
* `timeout_member_connect` - See Argument Reference above.
* `timeout_member_data` - See Argument Reference above.
* `timeout_tcp_inspect` - See Argument Reference above.
* `default_tls_container_ref` - See Argument Reference above.
* `sni_container_refs` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `insert_headers` - See Argument Reference above.
* `allowed_cidrs` - See Argument Reference above.

## Import

Load Balancer Listener can be imported using the Listener ID, e.g.:

```
$ terraform import vkcs_lb_listener.listener_1 b67ce64e-8b26-405d-afeb-4a078901f15a
```
