resource "vkcs_networking_floatingip" "associated_fip" {
  pool    = "internet"
  port_id = vkcs_networking_port.persistent_etcd.id
}
