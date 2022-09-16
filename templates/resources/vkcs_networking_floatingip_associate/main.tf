resource "vkcs_networking_port" "port_1" {
  network_id = "a5bbd213-e1d3-49b6-aed1-9df60ea94b9a"
}

resource "vkcs_networking_floatingip_associate" "fip_1" {
  floating_ip = "1.2.3.4"
  port_id     = "${vkcs_networking_port.port_1.id}"
}
