resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_port" "port_1" {
  name           = "port_1"
  network_id     = "${vkcs_networking_network.network_1.id}"
  admin_state_up = "true"
}
