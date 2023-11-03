data "vkcs_networking_port" "system_port" {
  fixed_ip = "10.0.0.10"
}

data "vkcs_networking_secgroup" "secgroup" {
  name = "secgroup"
}

resource "vkcs_networking_port_secgroup_associate" "port_1" {
  port_id = "${data.vkcs_networking_port.system_port.id}"
  security_group_ids = [
    "${data.vkcs_networking_secgroup.secgroup.id}",
  ]
}
