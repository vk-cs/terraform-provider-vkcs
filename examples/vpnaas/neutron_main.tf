resource "vkcs_vpnaas_ike_policy" "neutron" {
  name        = "neutron-ike-tf-example"
  ike_version = "v2"
  lifetime {
    units = "seconds"
    value = 3600
  }
  auth_algorithm          = local.auth_algorithm
  encryption_algorithm    = local.encryption_algorithm
  phase1_negotiation_mode = local.phase1_negotiation_mode
  sdn                     = local.sdn_neutron
}

resource "vkcs_vpnaas_ipsec_policy" "neutron" {
  name = "neutron-ipsec-tf-example"
  lifetime {
    units = "seconds"
    value = 3600
  }
  auth_algorithm       = local.auth_algorithm
  encryption_algorithm = local.encryption_algorithm
  pfs                  = local.pfs
  sdn                  = local.sdn_neutron
}

# Local endpoint
resource "vkcs_vpnaas_endpoint_group" "local_neutron" {
  name = "neutron-local-endpoint-tf-example"
  type = "subnet"
  endpoints = [
    vkcs_networking_subnet.neutron.id
  ]
  sdn = local.sdn_neutron

  # Wait for subnet is connected to the router
  depends_on = [
    vkcs_networking_router_interface.neutron
  ]
}

# Remote endpoint
resource "vkcs_vpnaas_endpoint_group" "remote_neutron" {
  name = "neutron-remote-endpoint-tf-example"
  type = "cidr"
  endpoints = [
    local.sprut_cidr
  ]
  sdn = local.sdn_neutron
}

resource "vkcs_vpnaas_service" "neutron" {
  name      = "neutron-vpn-tf-example"
  router_id = vkcs_networking_router.neutron.id
  sdn       = local.sdn_neutron
}

resource "vkcs_vpnaas_site_connection" "connection_neutron" {
  name              = "connection-neutron-tf-example"
  ikepolicy_id      = vkcs_vpnaas_ike_policy.neutron.id
  ipsecpolicy_id    = vkcs_vpnaas_ipsec_policy.neutron.id
  vpnservice_id     = vkcs_vpnaas_service.neutron.id
  psk               = local.psk_key
  peer_address      = vkcs_dc_interface.internet_sprut.ip_address
  peer_id           = vkcs_dc_interface.internet_sprut.ip_address
  local_ep_group_id = vkcs_vpnaas_endpoint_group.local_neutron.id
  peer_ep_group_id  = vkcs_vpnaas_endpoint_group.remote_neutron.id
  sdn               = local.sdn_neutron

  dpd {
    action   = "restart"
    timeout  = 60
    interval = 30
  }
}

# Static route in subnet to local network is very recommended.
# See https://cloud.vk.com/docs/networks/vnet/how-to-guides/vpn-tunnel#4_add_static_routes
# Find IP address of SNAT port of the router
data "vkcs_networking_port" "snat_port_neutron" {
  network_id   = vkcs_networking_network.neutron.id
  device_owner = "network:router_centralized_snat"

  # Wait for subnet is connected to the router
  depends_on = [
    vkcs_networking_router_interface.neutron
  ]
}

resource "vkcs_networking_subnet_route" "static_neutron" {
  subnet_id        = vkcs_networking_subnet.neutron.id
  destination_cidr = vkcs_networking_subnet.sprut.cidr
  next_hop         = data.vkcs_networking_port.snat_port_neutron.all_fixed_ips[0]
}
