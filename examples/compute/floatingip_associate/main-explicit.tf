resource "vkcs_networking_floatingip" "fip_explicit" {
  pool = "internet"
}

resource "vkcs_compute_floatingip_associate" "fip_explicit" {
  floating_ip = vkcs_networking_floatingip.fip_explicit.address
  instance_id = vkcs_compute_instance.multiple_networks.id
  fixed_ip    = vkcs_compute_instance.multiple_networks.network[1].fixed_ip_v4
}
