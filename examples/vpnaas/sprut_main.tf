resource "vkcs_vpnaas_ike_policy" "sprut" {
  name        = "sprut-ike-tf-example"
  ike_version = "v2"
  lifetime {
    units = "seconds"
    value = 3600
  }
  auth_algorithm          = local.auth_algorithm
  encryption_algorithm    = local.encryption_algorithm
  phase1_negotiation_mode = local.phase1_negotiation_mode
  sdn                     = local.sdn_sprut
}

resource "vkcs_vpnaas_ipsec_policy" "sprut" {
  name = "sprut-ipsec-tf-example"
  lifetime {
    units = "seconds"
    value = 3600
  }
  auth_algorithm       = local.auth_algorithm
  encryption_algorithm = local.encryption_algorithm
  pfs                  = local.pfs
  sdn                  = local.sdn_sprut
}

# Local sprut endpoint
resource "vkcs_vpnaas_endpoint_group" "local_sprut" {
  name = "sprut-local-endpoint-tf-example"
  type = "cidr"
  endpoints = [
    local.sprut_cidr
  ]
  sdn = local.sdn_sprut
}

# Remote sprut endpoint
resource "vkcs_vpnaas_endpoint_group" "remote_sprut" {
  name = "sprut-remote-endpoint-tf-example"
  type = "cidr"
  endpoints = [
    local.neutron_cidr
  ]
  sdn = local.sdn_sprut
}

resource "vkcs_vpnaas_service" "sprut" {
  name      = "sprut-vpn-tf-example"
  router_id = vkcs_dc_router.sprut.id
  sdn       = local.sdn_sprut

  depends_on = [
    vkcs_dc_interface.internet_sprut
  ]
}

resource "vkcs_vpnaas_site_connection" "connection_sprut" {
  name              = "connection-sprut-tf-example"
  ikepolicy_id      = vkcs_vpnaas_ike_policy.sprut.id
  ipsecpolicy_id    = vkcs_vpnaas_ipsec_policy.sprut.id
  vpnservice_id     = vkcs_vpnaas_service.sprut.id
  psk               = local.psk_key
  peer_address      = vkcs_networking_router.neutron.external_fixed_ips[0].ip_address
  peer_id           = vkcs_networking_router.neutron.external_fixed_ips[0].ip_address
  local_ep_group_id = vkcs_vpnaas_endpoint_group.local_sprut.id
  peer_ep_group_id  = vkcs_vpnaas_endpoint_group.remote_sprut.id
  sdn               = local.sdn_sprut

  dpd {
    action   = "restart"
    timeout  = 60
    interval = 30
  }
}

# We have to specify a path to DC router for VPN traffic
# since it is not a gateway in the network.
resource "vkcs_networking_subnet_route" "static_sprut" {
  subnet_id        = vkcs_networking_subnet.sprut.id
  destination_cidr = local.neutron_cidr
  next_hop         = vkcs_dc_interface.subnet_sprut.ip_address
}
