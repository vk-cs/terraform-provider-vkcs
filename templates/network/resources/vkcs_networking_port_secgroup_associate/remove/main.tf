data "vkcs_networking_port" "system_port" {
  fixed_ip = "10.0.0.10"
}

resource "vkcs_networking_port_secgroup_associate" "port_1" {
  port_id            = "${data.vkcs_networking_port.system_port.id}"
  enforce            = "true"
  security_group_ids = []
}
