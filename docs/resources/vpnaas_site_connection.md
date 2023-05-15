---
subcategory: "VPN"
layout: "vkcs"
page_title: "vkcs: vkcs_vpnaas_site_connection"
description: |-
  Manages a IPSec site connection resource within VKCS.
---

# vkcs_vpnaas_site_connection

Manages a IPSec site connection resource within VKCS.

## Example Usage
```terraform
resource "vkcs_vpnaas_service" "service" {
  router_id = "${vkcs_networking_router.router.id}"
}

resource "vkcs_vpnaas_ipsec_policy" "policy_1" {
	name = "ipsec-policy"
}

resource "vkcs_vpnaas_ike_policy" "policy_2" {
  name = "ike-policy"
}

resource "vkcs_vpnaas_endpoint_group" "group_1" {
	type = "cidr"
	endpoints = ["10.0.0.24/24", "10.0.0.25/24"]
}
resource "vkcs_vpnaas_endpoint_group" "group_2" {
	type = "subnet"
	endpoints = [ "${vkcs_networking_subnet.subnet.id}" ]
}

resource "vkcs_vpnaas_site_connection" "connection" {
	name = "connection"
	ikepolicy_id = "${vkcs_vpnaas_ike_policy.policy_2.id}"
	ipsecpolicy_id = "${vkcs_vpnaas_ipsec_policy.policy_1.id}"
	vpnservice_id = "${vkcs_vpnaas_service.service.id}"
	psk = "secret"
	peer_address = "192.168.10.1"
	peer_id = "192.168.10.1"
	local_ep_group_id = "${vkcs_vpnaas_endpoint_group.group_2.id}"
	peer_ep_group_id = "${vkcs_vpnaas_endpoint_group.group_1.id}"
	dpd {
		action   = "restart"
		timeout  = 42
		interval = 21
	}
	depends_on = ["vkcs_networking_router_interface.router_interface"]
}
```
## Argument Reference
- `ikepolicy_id` **required** *string* &rarr;  The ID of the IKE policy. Changing this creates a new connection.

- `ipsecpolicy_id` **required** *string* &rarr;  The ID of the IPsec policy. Changing this creates a new connection.

- `peer_address` **required** *string* &rarr;  The peer gateway public IPv4 or IPv6 address or FQDN.

- `peer_id` **required** *string* &rarr;  The peer router identity for authentication. A valid value is an IPv4 address, IPv6 address, e-mail address, key ID, or FQDN. Typically, this value matches the peer_address value. Changing this updates the existing policy.

- `psk` **required** *string* &rarr;  The pre-shared key. A valid value is any string.

- `vpnservice_id` **required** *string* &rarr;  The ID of the VPN service. Changing this creates a new connection.

- `admin_state_up` optional *boolean* &rarr;  The administrative state of the resource. Can either be up(true) or down(false). Changing this updates the administrative state of the existing connection.

- `description` optional *string* &rarr;  The human-readable description for the connection. Changing this updates the description of the existing connection.

- `dpd` optional &rarr;  A dictionary with dead peer detection (DPD) protocol controls.
  - `action` optional *string* &rarr;  The dead peer detection (DPD) action. A valid value is clear, hold, restart, disabled, or restart-by-peer. Default value is hold.

  - `interval` optional *number* &rarr;  The dead peer detection (DPD) interval, in seconds. A valid value is a positive integer. Default is 30.

  - `timeout` optional *number* &rarr;  The dead peer detection (DPD) timeout in seconds. A valid value is a positive integer that is greater than the DPD interval value. Default is 120.

- `initiator` optional *string* &rarr;  A valid value is response-only or bi-directional. Default is bi-directional.

- `local_ep_group_id` optional *string* &rarr;  The ID for the endpoint group that contains private subnets for the local side of the connection. You must specify this parameter with the peer_ep_group_id parameter unless in backward- compatible mode where peer_cidrs is provided with a subnet_id for the VPN service. Changing this updates the existing connection.

- `local_id` optional *string* &rarr;  An ID to be used instead of the external IP address for a virtual router used in traffic between instances on different networks in east-west traffic. Most often, local ID would be domain name, email address, etc. If this is not configured then the external IP address will be used as the ID.

- `mtu` optional *number* &rarr;  The maximum transmission unit (MTU) value to address fragmentation. Minimum value is 68 for IPv4, and 1280 for IPv6.

- `name` optional *string* &rarr;  The name of the connection. Changing this updates the name of the existing connection.

- `peer_cidrs` optional *string* &rarr;  Unique list of valid peer private CIDRs in the form < net_address > / < prefix >.

- `peer_ep_group_id` optional *string* &rarr;  The ID for the endpoint group that contains private CIDRs in the form < net_address > / < prefix > for the peer side of the connection. You must specify this parameter with the local_ep_group_id parameter unless in backward-compatible mode where peer_cidrs is provided with a subnet_id for the VPN service.

- `region` optional *string* &rarr;  The region in which to obtain the Networking client. A Networking client is needed to create an IPSec site connection. If omitted, the `region` argument of the provider is used. Changing this creates a new site connection.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Services can be imported using the `id`, e.g.

```shell
terraform import vkcs_vpnaas_site_connection.conn_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
