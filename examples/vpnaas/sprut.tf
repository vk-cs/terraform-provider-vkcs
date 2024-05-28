# Create sprut network
resource "vkcs_networking_network" "sprut" {
  name = "sprut-network-tf-example"
  sdn  = local.sdn_sprut
}

resource "vkcs_networking_subnet" "sprut" {
  name        = "sprut-subnet-tf-example"
  network_id  = vkcs_networking_network.sprut.id
  cidr        = "172.16.0.0/29"
  enable_dhcp = true
  allocation_pool {
    start = "172.16.0.2"
    end   = "172.16.0.6"
  }
  gateway_ip = "172.16.0.1"
  sdn        = local.sdn_sprut
}

# Get external network with Internet access
data "vkcs_networking_network" "internet_sprut" {
  name = "internet"
  sdn  = local.sdn_sprut
}

# Create a router to connect networks
resource "vkcs_dc_router" "sprut" {
  availability_zone = "GZ1"
  flavor            = "standard"
  name              = "dc-router-sprut-tf-example"
  description       = "dc_router in sprut"
}

# Connect internet to the router
resource "vkcs_dc_interface" "internet_sprut" {
  name         = "dc-interface-for-internet-sprut-tf-example"
  description  = "dc_interface for connecting dc_router to the internet"
  dc_router_id = vkcs_dc_router.sprut.id
  network_id   = data.vkcs_networking_network.internet_sprut.id
}

# Connect networks to the router
resource "vkcs_dc_interface" "subnet_sprut" {
  name         = "dc-interface-for-subnet-sprut-tf-example"
  description  = "dc_interface for connecting dc_router to the network and subnet"
  dc_router_id = vkcs_dc_router.sprut.id
  network_id   = vkcs_networking_network.sprut.id
  subnet_id    = vkcs_networking_subnet.sprut.id
}

############### Configure VPN ################
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
    vkcs_networking_subnet.sprut.cidr
  ]
  sdn = local.sdn_sprut
}

# Remote sprut endpoint
resource "vkcs_vpnaas_endpoint_group" "remote_sprut" {
  name = "sprut-remote-endpoint-tf-example"
  type = "cidr"
  endpoints = [
    vkcs_networking_subnet.neutron.cidr
  ]
  sdn = local.sdn_sprut
}

resource "vkcs_vpnaas_service" "sprut" {
  name      = "sprut-vpn-tf-example"
  router_id = vkcs_dc_router.sprut.id
  sdn       = local.sdn_sprut

  depends_on = [
    vkcs_dc_interface.subnet_sprut,
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

  depends_on = [
    vkcs_dc_interface.subnet_sprut,
    vkcs_dc_interface.internet_sprut,
    vkcs_vpnaas_ike_policy.neutron,
    vkcs_vpnaas_ipsec_policy.neutron,
    vkcs_vpnaas_service.neutron
  ]
}

############# Add static route ##############
resource "vkcs_networking_subnet_route" "static_sprut" {
  subnet_id        = vkcs_networking_subnet.sprut.id
  destination_cidr = vkcs_networking_subnet.neutron.cidr
  next_hop         = vkcs_dc_interface.subnet_sprut.ip_address
}

################## Verify ##################
# Create vm in sprut net with default, ssh and ping security groups.
# This is only necessary to verify that the vpn tunnel is configured correctly.
resource "vkcs_networking_secgroup" "ping_sprut" {
  name = "ping-sprut-tf-example"
  sdn  = local.sdn_sprut
}

resource "vkcs_networking_secgroup_rule" "ping_rule_sprut" {
  description       = "Ping rule in sprut"
  security_group_id = vkcs_networking_secgroup.ping_sprut.id
  direction         = "ingress"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  sdn               = local.sdn_sprut
}

resource "vkcs_networking_secgroup" "ssh_sprut" {
  name = "ssh-sprut-tf-example"
  sdn  = local.sdn_sprut
}

resource "vkcs_networking_secgroup_rule" "ssh_rule_sprut" {
  description       = "SSH rule in sprut"
  security_group_id = vkcs_networking_secgroup.ssh_sprut.id
  direction         = "ingress"
  protocol          = "tcp"
  # Specify SSH port
  port_range_max = 22
  port_range_min = 22
  # Allow access from any sources
  remote_ip_prefix = "0.0.0.0/0"
}

data "vkcs_networking_secgroup" "default_sprut" {
  name = "default"
  sdn  = local.sdn_sprut
}

# Create vm in sprut net
resource "vkcs_compute_instance" "vm_sprut" {
  name              = "sprut-vm-tf-example"
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
    uuid = vkcs_networking_network.sprut.id
  }

  security_group_ids = [
    data.vkcs_networking_secgroup.default_sprut.id,
    vkcs_networking_secgroup.ssh_sprut.id,
    vkcs_networking_secgroup.ping_sprut.id
  ]

  depends_on = [
    vkcs_dc_interface.subnet_sprut,
    vkcs_dc_interface.internet_sprut
  ]
}
