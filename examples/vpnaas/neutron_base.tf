# Create neutron network
resource "vkcs_networking_network" "neutron" {
  name = "neutron-network-tf-example"
  sdn  = local.sdn_neutron
}

resource "vkcs_networking_subnet" "neutron" {
  name        = "neutron-subnet-tf-example"
  network_id  = vkcs_networking_network.neutron.id
  cidr        = local.neutron_cidr
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
}
