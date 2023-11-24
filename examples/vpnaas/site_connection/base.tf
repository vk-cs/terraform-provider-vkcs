data "vkcs_networking_network" "extnet" {
  name = "ext-net"
  sdn = "neutron"
}

resource "vkcs_networking_network" "network" {
  name = "vpnaas_network"
  sdn = "neutron"
}

resource "vkcs_networking_subnet" "subnet" {
 	network_id = "${vkcs_networking_network.network.id}"
 	cidr = "192.168.199.0/24"
  sdn = "neutron"
}

resource "vkcs_networking_router" "router" {
  name = "router"
  external_network_id = data.vkcs_networking_network.extnet.id
  sdn = "neutron"
}

resource "vkcs_networking_router_interface" "router_interface" {
 	router_id = "${vkcs_networking_router.router.id}"
 	subnet_id = "${vkcs_networking_subnet.subnet.id}"
  sdn = "neutron"
}
