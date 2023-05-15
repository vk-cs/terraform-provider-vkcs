resource "vkcs_networking_router" "router_1" {
  name           = "router_1"
  admin_state_up = "true"
}

resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  network_id = "${vkcs_networking_network.network_1.id}"
  cidr       = "192.168.199.0/24"
}

resource "vkcs_networking_router_interface" "int_1" {
  router_id = "${vkcs_networking_router.router_1.id}"
  subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}

resource "vkcs_networking_router_route" "router_route_1" {
  depends_on       = ["vkcs_networking_router_interface.int_1"]
  router_id        = "${vkcs_networking_router.router_1.id}"
  destination_cidr = "10.0.1.0/24"
  next_hop         = "192.168.199.254"
}
