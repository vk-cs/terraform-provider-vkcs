resource "vkcs_vpnaas_service" "service" {
	router_id = "${vkcs_networking_router.router.id}"
	sdn = "neutron"
}

resource "vkcs_vpnaas_ipsec_policy" "policy_1" {
	name = "ipsec-policy"
	sdn = "neutron"
}

resource "vkcs_vpnaas_ike_policy" "policy_2" {
	name = "ike-policy"
	sdn = "neutron"
}

resource "vkcs_vpnaas_endpoint_group" "group_1" {
	type = "cidr"
	endpoints = ["10.0.0.24/24", "10.0.0.25/24"]
	sdn = "neutron"
}
resource "vkcs_vpnaas_endpoint_group" "group_2" {
	type = "subnet"
	endpoints = [ "${vkcs_networking_subnet.subnet.id}" ]
	sdn = "neutron"
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
	sdn = "neutron"
	depends_on = ["vkcs_networking_router_interface.router_interface"]
}
