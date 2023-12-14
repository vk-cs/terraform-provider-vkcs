resource "vkcs_compute_interface_attach" "db" {
  instance_id = vkcs_compute_instance.basic.id
  network_id  = vkcs_networking_network.db.id
}
