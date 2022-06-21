data "vkcs_networking_network" "extnet" {
  name = "ext-net"
}

resource "vkcs_networking_network" "network" {
  name = "vpnaas_network"
}

resource "vkcs_networking_subnet" "subnet" {
 	network_id = "${vkcs_networking_network.network.id}"
 	cidr = "192.168.199.0/24"
 	ip_version = 4
}

resource "vkcs_networking_router" "router" {
  name = "router"
  external_network_id = data.vkcs_networking_network.extnet.id
}

resource "vkcs_networking_router_interface" "router_interface" {
 	router_id = "${vkcs_networking_router.router.id}"
 	subnet_id = "${vkcs_networking_subnet.subnet.id}"
}
