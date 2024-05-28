# Create neutron network
resource "vkcs_networking_network" "neutron" {
  name = "neutron-network-tf-example"
  sdn  = local.sdn_neutron
}

resource "vkcs_networking_subnet" "neutron" {
  name        = "neutron-subnet-tf-example"
  network_id  = vkcs_networking_network.neutron.id
  cidr        = "10.0.0.0/29"
  enable_dhcp = true
  allocation_pool {
    start = "10.0.0.2"
    end   = "10.0.0.6"
  }
  gateway_ip = "10.0.0.1"
  sdn        = local.sdn_neutron
}

# Get external network with Internet access
data "vkcs_networking_network" "extnet_neutron" {
  name = "ext-net"
  sdn  = local.sdn_neutron
}

# Create a router with connection to Internet
resource "vkcs_networking_router" "neutron" {
  name                = "router-neutron-tf-example"
  external_network_id = data.vkcs_networking_network.extnet_neutron.id
  sdn                 = local.sdn_neutron
}

# Connect networks to the router
resource "vkcs_networking_router_interface" "neutron" {
  router_id = vkcs_networking_router.neutron.id
  subnet_id = vkcs_networking_subnet.neutron.id
  sdn       = local.sdn_neutron
}

############### Configure VPN ################
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

# Local neutron endpoint
resource "vkcs_vpnaas_endpoint_group" "local_neutron" {
  name      = "neutron-local-endpoint-tf-example"
  type      = "subnet"
  endpoints = [vkcs_networking_subnet.neutron.id]
  sdn       = local.sdn_neutron
}

# Remote neutron endpoint
resource "vkcs_vpnaas_endpoint_group" "remote_neutron" {
  name = "neutron-remote-endpoint-tf-example"
  type = "cidr"
  endpoints = [
    vkcs_networking_subnet.sprut.cidr
  ]
  sdn = local.sdn_neutron
}

resource "vkcs_vpnaas_service" "neutron" {
  name       = "neutron-vpn-tf-example"
  router_id  = vkcs_networking_router.neutron.id
  sdn        = local.sdn_neutron
  depends_on = [vkcs_networking_router_interface.neutron]
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

  depends_on = [
    vkcs_networking_router_interface.neutron,
    vkcs_vpnaas_ike_policy.sprut,
    vkcs_vpnaas_ipsec_policy.sprut,
    vkcs_vpnaas_service.sprut
  ]
}

############# Add static route ##############
# Find IP address of SNAT port of the router
data "vkcs_networking_port" "snat_port_neutron" {
  network_id   = vkcs_networking_network.neutron.id
  device_owner = "network:router_centralized_snat"
  depends_on = [
    vkcs_networking_network.neutron,
    vkcs_networking_router_interface.neutron
  ]
}

resource "vkcs_networking_subnet_route" "static_neutron" {
  subnet_id        = vkcs_networking_subnet.neutron.id
  destination_cidr = vkcs_networking_subnet.sprut.cidr
  next_hop         = data.vkcs_networking_port.snat_port_neutron.all_fixed_ips[0]
}

################## Verify ##################
# Create vm in neutron net with default, ssh and ping security groups.
# This is only necessary to verify that the vpn tunnel is configured correctly.
resource "vkcs_networking_secgroup" "ping_neutron" {
  name = "ping-neutron-tf-example"
  sdn  = local.sdn_neutron
}

resource "vkcs_networking_secgroup_rule" "ping_rule_neutron" {
  description       = "Ping rule in neutron"
  security_group_id = vkcs_networking_secgroup.ping_neutron.id
  direction         = "ingress"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  sdn               = local.sdn_neutron
}

resource "vkcs_networking_secgroup" "ssh_neutron" {
  name = "ssh-neutron-tf-example"
  sdn  = local.sdn_neutron
}

resource "vkcs_networking_secgroup_rule" "ssh_rule_neutron" {
  description       = "SSH rule in neutron"
  security_group_id = vkcs_networking_secgroup.ssh_neutron.id
  direction         = "ingress"
  protocol          = "tcp"
  # Specify SSH port
  port_range_max = 22
  port_range_min = 22
  # Allow access from any sources
  remote_ip_prefix = "0.0.0.0/0"
}

data "vkcs_networking_secgroup" "default_neutron" {
  name = "default"
  sdn  = local.sdn_neutron
}

# Create vm in neutron net
resource "vkcs_compute_instance" "vm_neutron" {
  name              = "neutron-vm-tf-example"
  availability_zone = "ME1"
  flavor_name       = "STD3-1-2"
  # key_pair        = <name_of_your_key_pair> # need for connect to vm by ssh
  block_device {
    source_type           = "image"
    uuid                  = data.vkcs_images_image.image.id
    destination_type      = "volume"
    volume_size           = 11
    delete_on_termination = true
  }

  network {
    uuid = vkcs_networking_network.neutron.id
  }

  security_group_ids = [
    data.vkcs_networking_secgroup.default_neutron.id,
    vkcs_networking_secgroup.ssh_neutron.id,
    vkcs_networking_secgroup.ping_neutron.id
  ]

  depends_on = [
    vkcs_networking_router_interface.neutron
  ]
}
