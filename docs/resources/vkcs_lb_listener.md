---
layout: "vkcs"
page_title: "vkcs: vkcs_lb_listener"
description: |-
  Manages a listener resource within VKCS.
---

# vkcs_lb_listener

Manages a listener resource within VKCS.

## Example Usage
```terraform
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
- `loadbalancer_id` **String** (***Required***) The load balancer on which to provision this Listener. Changing this creates a new Listener.

- `protocol` **String** (***Required***) The protocol - can either be TCP, HTTP, HTTPS, TERMINATED_HTTPS, UDP. Changing this creates a new Listener.

- `protocol_port` **Number** (***Required***) The port on which to listen for client traffic. Changing this creates a new Listener.

- `admin_state_up` **Boolean** (*Optional*) The administrative state of the Listener. A valid value is true (UP) or false (DOWN).

- `allowed_cidrs` **String** (*Optional*) A list of CIDR blocks that are permitted to connect to this listener, denying all other source addresses. If not present, defaults to allow all.

- `connection_limit` **Number** (*Optional*) The maximum number of connections allowed for the Listener.

- `default_pool_id` **String** (*Optional*) The ID of the default pool with which the Listener is associated.

- `default_tls_container_ref` **String** (*Optional*) A reference to a Keymanager Secrets container which stores TLS information. This is required if the protocol is `TERMINATED_HTTPS`.

- `description` **String** (*Optional*) Human-readable description for the Listener.

- `insert_headers` <strong>Map of </strong>**String** (*Optional*) The list of key value pairs representing headers to insert into the request before it is sent to the backend members. Changing this updates the headers of the existing listener.

- `name` **String** (*Optional*) Human-readable name for the Listener. Does not have to be unique.

- `region` **String** (*Optional*) The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new Listener.

- `sni_container_refs` **String** (*Optional*) A list of references to Keymanager Secrets containers which store SNI information.

- `timeout_client_data` **Number** (*Optional*) The client inactivity timeout in milliseconds.

- `timeout_member_connect` **Number** (*Optional*) The member connection timeout in milliseconds.

- `timeout_member_data` **Number** (*Optional*) The member inactivity timeout in milliseconds.

- `timeout_tcp_inspect` **Number** (*Optional*) The time in milliseconds, to wait for additional TCP packets for content inspection.


## Attributes Reference
- `loadbalancer_id` **String** See Argument Reference above.

- `protocol` **String** See Argument Reference above.

- `protocol_port` **Number** See Argument Reference above.

- `admin_state_up` **Boolean** See Argument Reference above.

- `allowed_cidrs` **String** See Argument Reference above.

- `connection_limit` **Number** See Argument Reference above.

- `default_pool_id` **String** See Argument Reference above.

- `default_tls_container_ref` **String** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `insert_headers` <strong>Map of </strong>**String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `sni_container_refs` **String** See Argument Reference above.

- `timeout_client_data` **Number** See Argument Reference above.

- `timeout_member_connect` **Number** See Argument Reference above.

- `timeout_member_data` **Number** See Argument Reference above.

- `timeout_tcp_inspect` **Number** See Argument Reference above.

- `id` **String** ID of the resource.



## Import

Load Balancer Listener can be imported using the Listener ID, e.g.:

```shell
terraform import vkcs_lb_listener.listener_1 b67ce64e-8b26-405d-afeb-4a078901f15a
```
