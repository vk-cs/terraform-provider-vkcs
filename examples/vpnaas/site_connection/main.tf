resource "vkcs_vpnaas_site_connection" "connection" {
  name              = "connection"
  ikepolicy_id      = vkcs_vpnaas_ike_policy.data_center.id
  ipsecpolicy_id    = vkcs_vpnaas_ipsec_policy.data_center.id
  vpnservice_id     = vkcs_vpnaas_service.vpn_to_datacenter.id
  psk               = "secret"
  peer_address      = "192.168.10.1"
  peer_id           = "192.168.10.1"
  local_ep_group_id = vkcs_vpnaas_endpoint_group.subnet_hosts.id
  peer_ep_group_id  = vkcs_vpnaas_endpoint_group.allowed_hosts.id
  dpd {
    action   = "restart"
    timeout  = 42
    interval = 21
  }
  depends_on = [vkcs_dc_interface.internet]
}
