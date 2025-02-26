resource "vkcs_networking_floatingip" "fip_basic" {
  pool = "internet"
}

resource "vkcs_compute_floatingip_associate" "fip_basic" {
  floating_ip = vkcs_networking_floatingip.fip_basic.address
  instance_id = vkcs_compute_instance.basic.id
}
