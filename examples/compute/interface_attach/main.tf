resource "vkcs_compute_interface_attach" "etcd" {
  instance_id = vkcs_compute_instance.basic.id
  port_id     = vkcs_networking_port.persistent_etcd.id
}
